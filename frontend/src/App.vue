<script setup lang="ts">
import { ref } from 'vue';
import HeaderToolbar from './components/HeaderToolbar.vue';
import MainLayout from './components/MainLayout.vue';
import type { ViewMode } from './types/viewMode';
import { setProjectRoot } from './utils/markdownUtils';

// 状態管理
const rootPath = ref<string>('');
const isLoading = ref(false);
const viewMode = ref<ViewMode>('split');
const hasActiveFile = ref(false);

// 表示モードを変更する
const changeViewMode = (mode: ViewMode): void => {
  viewMode.value = mode;
};

// フォルダが変更された時の処理
const handleFolderChanged = (path: string): void => {
  rootPath.value = path;
  setProjectRoot(rootPath.value);
};

// ファイルがアクティブになった時の処理
const handleFileActiveChanged = (isActive: boolean): void => {
  hasActiveFile.value = isActive;
};

// FileExplorerのフォルダを開くボタンのハンドラー（HeaderToolbarからの呼び出し用）
const handleOpenFolder = (): void => {
  // FileExplorerのopenFolderメソッドを呼び出すためのイベント
  // 実際の処理はFileExplorerで行われる
};
</script>

<template>
  <div class="app-container">
    <!-- ヘッダーツールバー -->
    <HeaderToolbar
      :has-active-file="hasActiveFile"
      :view-mode="viewMode"
      @open-folder="handleOpenFolder"
      @change-view-mode="changeViewMode"
    />

    <!-- メインレイアウト -->
    <MainLayout
      :root-path="rootPath"
      :is-loading="isLoading"
      :view-mode="viewMode"
      @folder-changed="handleFolderChanged"
      @file-active-changed="handleFileActiveChanged"
    />
  </div>
</template>

<style>
/* リセットとベース */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  /* ライトモードのカラーテーマ */
  --bg-color: #ffffff;
  --sidebar-bg: #f5f5f5;
  --sidebar-header-bg: #e8e8e8;
  --text-color: #333333;
  --border-color: #dddddd;
  --accent-color: #0078d7;
  --hover-bg: #eaeaea;
  --header-bg: #f0f0f0;
  --header-border: #d1d1d1;
  --tab-bg: #f8f8f8;
  --tab-active-bg: #ffffff;
  --tab-hover-bg: #e8e8e8;
  --tab-border: #dddddd;
}

html,
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  height: 100%;
  color: var(--text-color);
  background-color: var(--bg-color);
}

#app {
  height: 100vh;
  width: 100vw;
}

/* アプリ全体のレイアウト */
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
  overflow: hidden;
}
</style>
