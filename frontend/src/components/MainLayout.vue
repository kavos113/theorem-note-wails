<script setup lang="ts">
import { ref, watch } from 'vue';
import FileExplorer from './FileExplorer.vue';
import MarkdownEditor from './MarkdownEditor.vue';
import TabBar, { type OpenFile } from './TabBar.vue';
import type { ViewMode } from '../types/viewMode';

interface Props {
  rootPath: string;
  isLoading: boolean;
  viewMode: ViewMode;
}

defineProps<Props>();

const emit = defineEmits<{
  'folder-changed': [path: string];
  'file-active-changed': [isActive: boolean];
}>();

const tabBarRef = ref<InstanceType<typeof TabBar> | null>(null);
const currentFile = ref<OpenFile | null>(null);
const selectedFilePath = ref<string | undefined>(undefined);

// アクティブファイルの変更を監視してイベントを発行
watch(currentFile, (newFile) => {
  emit('file-active-changed', !!newFile);
});

// ファイル選択時の処理
const handleFileSelect = async (filePath: string): Promise<void> => {
  if (tabBarRef.value) {
    await tabBarRef.value.openFileInTab(filePath);
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
          :selected-file-path="currentFile.path"
          :file-content="currentFile.content"
          :view-mode="viewMode"
          @update:file-content="handleContentUpdate"
          @file-saved="handleFileSaved"
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
