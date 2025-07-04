<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import FileTreeItem from './FileTreeItem.vue';
import {
  GetFileTree,
  GetLastOpened,
  GetNewDirectoryFileTree,
  CreateFile,
  CreateDirectory
} from '../../wailsjs/go/main/App';
import { backend } from '../../wailsjs/go/models';
import FileItem = backend.FileItem;

const props = defineProps<{
  rootPath?: string;
  selectedFile?: string;
}>();

const emit = defineEmits<{
  'select-file': [path: string, header?: string];
  'folder-changed': [path: string];
}>();

// 状態管理
const rootPath = ref<string>(props.rootPath || '');
const fileTree = ref<FileItem[]>([]);
const loading = ref<boolean>(false);
const error = ref<string | null>(null);
const selectedItem = ref<string | null>(null);
const expandedItems = ref<Set<string>>(new Set());

// フォルダを開く
const openFolder = async (): Promise<void> => {
  try {
    const result = await GetNewDirectoryFileTree();
    const newRootPath = await GetLastOpened();
    if (result && newRootPath) {
      rootPath.value = newRootPath;
      fileTree.value = result;
      emit('folder-changed', newRootPath);
    }
  } catch (err) {
    console.error('フォルダ選択でエラーが発生しました:', err);
    error.value = err instanceof Error ? err.message : '不明なエラーが発生しました';
  }
};

// 前回のディレクトリを読み込む
const loadLastDirectory = async (): Promise<void> => {
  try {
    loading.value = true;
    const lastDirectory = await GetLastOpened();
    if (lastDirectory) {
      const fileTreeResult = await GetFileTree(lastDirectory);
      rootPath.value = lastDirectory;
      fileTree.value = fileTreeResult;
      emit('folder-changed', lastDirectory);
    }
  } catch (err) {
    console.log('前回のディレクトリの読み込みに失敗:', err);
  } finally {
    loading.value = false;
  }
};

// ファイルツリーを読み込む
const loadFileTree = async (): Promise<void> => {
  if (!rootPath.value) return;
  try {
    loading.value = true;
    error.value = null;
    const result = await GetFileTree(rootPath.value);
    fileTree.value = result;
  } catch (err) {
    error.value = err instanceof Error ? err.message : '不明なエラーが発生しました';
    fileTree.value = [];
  } finally {
    loading.value = false;
  }
};

const getBasePath = (): string => {
  if (!selectedItem.value) {
    return rootPath.value;
  }
  // TODO: 現状、選択アイテムがファイルかフォルダかをフロントで判定できない。
  // そのため、選択アイテムがファイルの場合、その親ディレクトリを返すような処理ができない。
  // Go側でIsDirectoryを返すようにしているので、FileItemのツリーから探すことは可能だが複雑になる。
  // ここでは簡略化のため、選択アイテムのパスをそのままベースパスとする。
  // ユーザーがファイルを選択して「新規作成」すると、そのファイルと同じ階層に作成されることになる。
  return selectedItem.value;
};

const createFile = async () => {
  let fileName = prompt('新しいファイル名を入力してください');
  if (!fileName) return;

  if (!fileName.endsWith('.md')) {
    fileName += '.md';
  }

  const basePath = getBasePath();
  // Note: path.joinが使えないので手動で結合
  const newFilePath = basePath + '\\' + fileName;

  try {
    await CreateFile(newFilePath);
    await loadFileTree();
  } catch (err) {
    alert(`ファイル作成エラー: ${err}`);
  }
};

const createDirectory = async () => {
  const dirName = prompt('新しいフォルダ名を入力してください');
  if (!dirName) return;

  const basePath = getBasePath();
  const newDirPath = basePath + '\\' + dirName;

  try {
    await CreateDirectory(newDirPath);
    await loadFileTree();
  } catch (err) {
    alert(`フォルダ作成エラー: ${err}`);
  }
};

const handleExpandItem = (path: string) => {
  if (expandedItems.value.has(path)) {
    expandedItems.value.delete(path);
  } else {
    expandedItems.value.add(path);
  }
  console.log(expandedItems.value);
};

defineExpose({
  createFile,
  createDirectory
});

onMounted(() => {
  if (props.rootPath) {
    loadFileTree();
  } else {
    loadLastDirectory();
  }
});

watch(
  () => props.rootPath,
  (newRootPath) => {
    rootPath.value = newRootPath || '';
    selectedItem.value = null; // ルートが変わったら選択を解除
    if (newRootPath) {
      loadFileTree();
    } else {
      fileTree.value = [];
      loading.value = false;
    }
  }
);
</script>

<template>
  <div class="file-explorer">
    <div class="explorer-header">
      <h3>エクスプローラー</h3>
      <div>
        <button class="action-button" title="新しいファイル" @click="createFile">📄+</button>
        <button class="action-button" title="新しいフォルダ" @click="createDirectory">📁+</button>
        <button class="folder-button" title="フォルダを開く" @click="openFolder">📁</button>
      </div>
    </div>
    <div class="file-tree">
      <div v-if="loading" class="loading">読み込み中...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      <div v-else-if="!rootPath" class="no-folder">
        <p>フォルダが選択されていません</p>
        <button class="open-folder-btn" @click="openFolder">フォルダを開く</button>
      </div>
      <div v-else class="tree-container">
        <ul class="tree-root">
          <file-tree-item
            v-for="item in fileTree"
            :key="item.Path"
            :item="item"
            :selected-item="selectedItem"
            :expanded-items="expandedItems"
            @select-file="$emit('select-file', $event)"
            @select-item="selectedItem = $event"
            @expand-item="handleExpandItem"
          />
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.file-explorer {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--sidebar-bg, #f5f5f5);
  color: var(--text-color, #333333);
  user-select: none;
  overflow: hidden;
}

.explorer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  font-size: 13px;
  font-weight: 600;
  border-bottom: 1px solid var(--border-color, #dddddd);
}

.explorer-header h3 {
  margin: 0;
  font-size: inherit;
  font-weight: inherit;
}

.action-button,
.folder-button {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
  padding: 4px;
  border-radius: 3px;
  transition: background-color 0.2s;
  margin-left: 4px;
}

.action-button:hover,
.folder-button:hover {
  background-color: var(--hover-bg, #eaeaea);
}

.no-folder {
  padding: 20px;
  text-align: center;
  color: var(--text-color, #333333);
}

.no-folder p {
  margin-bottom: 15px;
  color: var(--text-muted, #666666);
}

.open-folder-btn {
  padding: 8px 16px;
  background-color: var(--accent-color, #0078d7);
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  transition: background-color 0.2s;
}

.open-folder-btn:hover {
  background-color: #0062a3;
}

.file-tree {
  flex: 1;
  overflow: auto;
  padding: 5px 0;
}

.tree-container {
  width: 100%;
}

.tree-root {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.loading,
.error {
  padding: 10px;
  color: var(--text-muted, #666666);
  font-style: italic;
}

.error {
  color: var(--error-text, #d32f2f);
}
</style>
