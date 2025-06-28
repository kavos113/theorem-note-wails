<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import FileTreeItem from './FileTreeItem.vue';
import { GetFileTree, GetLastOpened, GetNewDirectoryFileTree } from '../../wailsjs/go/main/App';
import { backend } from '../../wailsjs/go/models';
import FileItem = backend.FileItem;

const props = defineProps<{
  rootPath?: string;
  selectedFile?: string;
}>();

const emit = defineEmits<{
  'select-file': [path: string];
  'folder-changed': [path: string];
}>();

// 状態管理
const rootPath = ref<string>(props.rootPath || '');
const fileTree = ref<FileItem[]>([]);
const loading = ref<boolean>(false);
const error = ref<string | null>(null);

// フォルダを開く
const openFolder = async (): Promise<void> => {
  try {
    const fileTreeResult = await GetNewDirectoryFileTree();
    console.dir(fileTreeResult, { depth: null });
    if (fileTreeResult && fileTreeResult.length > 0) {
      // 最初のアイテムのパスからルートパスを抽出
      const firstPath = fileTreeResult[0].Path;
      const newRootPath = firstPath.split('\\').slice(0, -1).join('\\');
      rootPath.value = newRootPath;
      fileTree.value = fileTreeResult;
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
      // 前回のディレクトリのファイルツリーを取得
      const fileTreeResult = await GetFileTree(lastDirectory);
      if (fileTreeResult && fileTreeResult.length > 0) {
        rootPath.value = lastDirectory;
        fileTree.value = fileTreeResult;
        emit('folder-changed', lastDirectory);
      }
    }
  } catch (err) {
    console.log('前回のディレクトリの読み込みに失敗:', err);
    // エラーが発生した場合は何もしない（手動でフォルダを開く必要がある）
  } finally {
    loading.value = false;
  }
};

// ファイルツリーを読み込む
const loadFileTree = async (): Promise<void> => {
  try {
    loading.value = true;
    error.value = null;
    // ファイル一覧をElectronのメインプロセスから取得
    const result = await GetFileTree(rootPath.value);
    fileTree.value = result;
  } catch (err) {
    error.value = err instanceof Error ? err.message : '不明なエラーが発生しました';
    fileTree.value = [];
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  if (props.rootPath) {
    loadFileTree();
  } else {
    loadLastDirectory();
  }
});

// rootPathが変更された時にファイルツリーを更新
watch(
  () => props.rootPath,
  (newRootPath) => {
    rootPath.value = newRootPath || '';
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
      <button class="folder-button" title="フォルダを開く" @click="openFolder">📁</button>
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
            :selected-file="selectedFile"
            @select-file="$emit('select-file', $event)"
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

.folder-button {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
  padding: 4px;
  border-radius: 3px;
  transition: background-color 0.2s;
}

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
