<script setup lang="ts">
import type { ViewMode } from '../types/viewMode';

interface Props {
  hasActiveFile: boolean;
  viewMode: ViewMode;
}

defineProps<Props>();

defineEmits<{
  'open-folder': [];
  'change-view-mode': [mode: ViewMode];
  'open-settings': [];
}>();
</script>

<template>
  <div class="header-toolbar">
    <div class="toolbar-left">
      <button class="toolbar-button" @click="$emit('open-folder')">
        <span class="icon">📁</span>
        フォルダを開く
      </button>
    </div>
    <div class="toolbar-center">
      <h1 class="app-title">Theorem Note</h1>
    </div>
    <div class="toolbar-right">
      <div v-if="hasActiveFile" class="view-mode-buttons">
        <button
          class="toolbar-button"
          :class="{ active: viewMode === 'editor' }"
          @click="$emit('change-view-mode', 'editor')"
        >
          <span class="icon">📝</span>
          エディタのみ
        </button>
        <button
          class="toolbar-button"
          :class="{ active: viewMode === 'split' }"
          @click="$emit('change-view-mode', 'split')"
        >
          <span class="icon">⚡</span>
          分割表示
        </button>
        <button
          class="toolbar-button"
          :class="{ active: viewMode === 'preview' }"
          @click="$emit('change-view-mode', 'preview')"
        >
          <span class="icon">👁️</span>
          プレビューのみ
        </button>
      </div>
      <button class="toolbar-button" @click="$emit('open-settings')">
        <span class="icon">⚙️</span>
        設定
      </button>
    </div>
  </div>
</template>

<style scoped>
/* ヘッダーツールバー */
.header-toolbar {
  display: flex;
  align-items: center;
  height: 40px;
  background-color: var(--header-bg);
  border-bottom: 1px solid var(--header-border);
  padding: 0 10px;
  flex-shrink: 0;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-center {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-color);
  margin: 0;
}

.toolbar-button {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background-color: var(--accent-color);
  color: white;
  border: none;
  border-radius: 3px;
  font-size: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.toolbar-button:hover {
  background-color: #0062a3;
}

.toolbar-button .icon {
  font-size: 14px;
}

.toolbar-button.active {
  background-color: #005a9e;
}

.toolbar-button.active:hover {
  background-color: #004578;
}

.view-mode-buttons {
  display: flex;
  gap: 4px;
}

.view-mode-buttons .toolbar-button {
  min-width: 90px;
}
</style>
