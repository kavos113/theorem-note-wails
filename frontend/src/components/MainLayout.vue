<script setup lang="ts">
import { ref, watch } from 'vue';
import FileExplorer from './FileExplorer.vue';
import MarkdownEditor from './MarkdownEditor.vue';
import TabBar, { type OpenFile } from './TabBar.vue';
import type { ViewMode } from '../types/viewMode';
import { LoadSession, SaveSession } from '../../wailsjs/go/main/App';

interface Props {
  rootPath: string;
  isLoading: boolean;
  viewMode: ViewMode;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  'folder-changed': [path: string];
  'file-active-changed': [isActive: boolean];
}>();

const tabBarRef = ref<InstanceType<typeof TabBar> | null>(null);
const markdownEditorRef = ref<InstanceType<typeof MarkdownEditor> | null>(null);
const currentFile = ref<OpenFile | null>(null);
const selectedFilePath = ref<string | undefined>(undefined);

// アクティブファイルの変更を監視してイベントを発行
watch(currentFile, (newFile) => {
  emit('file-active-changed', !!newFile);
});

// ファイル選択時の処理
const handleFileSelect = async (filePath: string, header?: string): Promise<void> => {
  if (tabBarRef.value) {
    await tabBarRef.value.openFileInTab(filePath);

    if (header && markdownEditorRef.value) {
      // 少し待ってからスクロールを実行
      setTimeout(() => {
        markdownEditorRef.value?.scrollToHeader(header);
      }, 100); // DOMの更新を待つ
    }

    const filePaths = tabBarRef.value.openFiles.map((file) => file.path);
    try {
      await SaveSession(props.rootPath, filePaths);
    } catch (err) {
      console.error('セッションの保存に失敗しました:', err);
    }
  }
};

// タブバーからのイベント処理
const handleFileOpened = (file: OpenFile): void => {
  currentFile.value = file;
  selectedFilePath.value = file.path;
};

const handleFileClosed = (): void => {
  // タブが閉じられた後の現在のアクティブファイルを取得
  if (tabBarRef.value) {
    currentFile.value = tabBarRef.value.activeFile || null;
    selectedFilePath.value = currentFile.value?.path;
  }
};

const handleFileSwitched = (file: OpenFile): void => {
  currentFile.value = file;
  selectedFilePath.value = file.path;
};

const handleContentUpdated = (): void => {
  // タブバーで既に処理されているので、特に何もしない
};

// エディターからのコンテンツ更新
const handleContentUpdate = (newContent: string): void => {
  if (currentFile.value && tabBarRef.value) {
    tabBarRef.value.updateFileContent(currentFile.value.path, newContent);
  }
};

// ファイル保存後の処理
const handleFileSaved = (): void => {
  if (currentFile.value && tabBarRef.value) {
    tabBarRef.value.markFileAsSaved(currentFile.value.path);
  }
};

watch(
  () => props.rootPath,
  async (newPath) => {
    if (newPath && tabBarRef.value) {
      // 以前のセッションを読み込む
      try {
        const filePaths = await LoadSession(newPath);
        if (filePaths && filePaths.length > 0 && tabBarRef.value) {
          for (const path of filePaths) {
            await tabBarRef.value.openFileInTab(path);
          }
        }
      } catch (err) {
        console.error('セッションの読み込みに失敗しました:', err);
      }
    }
  }
);
</script>

<template>
  <div class="main-area">
    <!-- サイドバー（ファイルエクスプローラー） -->
    <div class="sidebar">
      <FileExplorer
        :root-path="rootPath"
        :selected-file="selectedFilePath"
        @select-file="handleFileSelect"
        @folder-changed="$emit('folder-changed', $event)"
      />
    </div>

    <!-- メインコンテンツ領域 -->
    <div class="main-content">
      <!-- タブバー -->
      <TabBar
        ref="tabBarRef"
        :is-loading="isLoading"
        @file-opened="handleFileOpened"
        @file-closed="handleFileClosed"
        @file-switched="handleFileSwitched"
        @content-updated="handleContentUpdated"
      />

      <div v-if="isLoading" class="loading">読み込み中...</div>
      <div v-else-if="!currentFile" class="welcome-screen">
        <h1>Theorem Note</h1>
        <p>
          左側のエクスプローラーからファイルを選択するか、「フォルダを開く」ボタンを押してプロジェクトフォルダを選択してください。
        </p>
      </div>
      <div v-else class="editor-container">
        <MarkdownEditor
          ref="markdownEditorRef"
          :selected-file-path="currentFile.path"
          :file-content="currentFile.content"
          :view-mode="viewMode"
          @update:file-content="handleContentUpdate"
          @file-saved="handleFileSaved"
          @select-file="handleFileSelect"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
/* メインエリア */
.main-area {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* サイドバー */
.sidebar {
  width: 250px;
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--sidebar-bg);
  border-right: 1px solid var(--border-color);
  overflow: hidden;
}

/* メインコンテンツ */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.welcome-screen {
  padding: 40px;
  text-align: center;
  color: var(--text-color);
}

.welcome-screen h1 {
  font-size: 2rem;
  margin-bottom: 20px;
}

.welcome-screen p {
  margin: 10px 0;
}

.editor-container {
  display: flex;
  flex-direction: column;
  height: calc(100% - 40px);
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  font-style: italic;
  color: var(--text-color);
}
</style>
