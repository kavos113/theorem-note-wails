<script setup lang="ts">
import { ref, computed } from 'vue';
import { ReadFile } from "../../wailsjs/go/main/App";

// タブで開いているファイルの型定義
export interface OpenFile {
  path: string;
  content: string;
  isModified: boolean;
  displayName: string;
}

interface Props {
  isLoading?: boolean;
}

defineProps<Props>();

const emit = defineEmits<{
  'file-opened': [file: OpenFile];
  'file-closed': [index: number];
  'file-switched': [file: OpenFile];
  'content-updated': [filePath: string, content: string];
}>();

// タブ管理状態
const openFiles = ref<OpenFile[]>([]);
const activeTabIndex = ref<number>(-1);

// 現在アクティブなファイルの計算プロパティ
const activeFile = computed((): OpenFile | undefined => {
  if (activeTabIndex.value >= 0 && activeTabIndex.value < openFiles.value.length) {
    return openFiles.value[activeTabIndex.value];
  }
  return undefined;
});

// ファイル名からタブ表示用の名前を取得
const getDisplayName = (filePath: string): string => {
  return filePath.split('\\').pop() || filePath.split('/').pop() || filePath;
};

// 新しいタブでファイルを開く
const openFileInTab = async (filePath: string): Promise<void> => {
  try {
    // 既に開いているファイルかチェック
    const existingIndex = openFiles.value.findIndex((file) => file.path === filePath);
    if (existingIndex !== -1) {
      // 既に開いている場合はそのタブをアクティブにする
      activeTabIndex.value = existingIndex;
      emit('file-switched', openFiles.value[existingIndex]);
      return;
    }

    // ファイルの内容を読み込む
    const content = await ReadFile(filePath);

    // 新しいタブを作成
    const newFile: OpenFile = {
      path: filePath,
      content,
      isModified: false,
      displayName: getDisplayName(filePath)
    };

    openFiles.value.push(newFile);
    activeTabIndex.value = openFiles.value.length - 1;
    emit('file-opened', newFile);
  } catch (err) {
    console.error('ファイル読み込みエラー:', err);
    // エラーの場合でもタブを作成（エラーメッセージを表示）
    const newFile: OpenFile = {
      path: filePath,
      content: '# エラー\nファイルを読み込めませんでした',
      isModified: false,
      displayName: getDisplayName(filePath)
    };
    openFiles.value.push(newFile);
    activeTabIndex.value = openFiles.value.length - 1;
    emit('file-opened', newFile);
  }
};

// タブを閉じる
const closeTab = (index: number): void => {
  if (index < 0 || index >= openFiles.value.length) return;

  openFiles.value.splice(index, 1);
  emit('file-closed', index);

  // アクティブタブのインデックスを調整
  if (openFiles.value.length === 0) {
    activeTabIndex.value = -1;
  } else if (index <= activeTabIndex.value) {
    if (activeTabIndex.value > 0) {
      activeTabIndex.value--;
    } else {
      activeTabIndex.value = 0;
    }
    // 新しいアクティブファイルを通知
    if (activeFile.value) {
      emit('file-switched', activeFile.value);
    }
  }
};

// タブを切り替える
const switchToTab = (index: number): void => {
  if (index >= 0 && index < openFiles.value.length) {
    activeTabIndex.value = index;
    emit('file-switched', openFiles.value[index]);
  }
};

// コンテンツ更新を処理
const updateFileContent = (filePath: string, newContent: string): void => {
  const fileIndex = openFiles.value.findIndex((f) => f.path === filePath);
  if (fileIndex !== -1) {
    openFiles.value[fileIndex].content = newContent;
    openFiles.value[fileIndex].isModified = true;
    emit('content-updated', filePath, newContent);
  }
};

// ファイル保存後に変更フラグをリセット
const markFileAsSaved = (filePath: string): void => {
  const fileIndex = openFiles.value.findIndex((f) => f.path === filePath);
  if (fileIndex !== -1) {
    openFiles.value[fileIndex].isModified = false;
  }
};

// 外部からの呼び出し用にメソッドを公開
defineExpose({
  openFileInTab,
  updateFileContent,
  markFileAsSaved,
  activeFile,
  openFiles,
  activeTabIndex
});
</script>

<template>
  <div v-if="openFiles.length > 0" class="tab-bar">
    <div
      v-for="(file, index) in openFiles"
      :key="file.path"
      class="tab"
      :class="{ active: index === activeTabIndex }"
      @click="switchToTab(index)"
    >
      <span class="tab-name">{{ file.displayName }}</span>
      <span v-if="file.isModified" class="modified-indicator">●</span>
      <button class="tab-close" @click.stop="closeTab(index)">×</button>
    </div>
  </div>
</template>

<style scoped>
/* タブバー */
.tab-bar {
  display: flex;
  background-color: var(--tab-bg);
  border-bottom: 1px solid var(--tab-border);
  overflow-x: auto;
  overflow-y: hidden;
  flex-shrink: 0;
  height: 40px;
  white-space: nowrap;
}

.tab-bar::-webkit-scrollbar {
  height: 3px;
}

.tab-bar::-webkit-scrollbar-track {
  background: var(--tab-bg);
}

.tab-bar::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

.tab-bar::-webkit-scrollbar-thumb:hover {
  background: var(--text-color);
}

.tab {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background-color: var(--tab-bg);
  border-right: 1px solid var(--tab-border);
  cursor: pointer;
  user-select: none;
  min-width: 120px;
  max-width: 200px;
  transition: background-color 0.2s;
  position: relative;
  flex-shrink: 0;
}

.tab:hover {
  background-color: var(--tab-hover-bg);
}

.tab.active {
  background-color: var(--tab-active-bg);
  border-bottom: 2px solid var(--accent-color);
  margin-bottom: -1px;
}

.tab-name {
  flex: 1;
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-right: 8px;
  color: var(--text-color);
}

.tab.active .tab-name {
  font-weight: 500;
}

.modified-indicator {
  color: var(--accent-color);
  font-size: 12px;
  margin-right: 4px;
  font-weight: bold;
}

.tab-close {
  background: none;
  border: none;
  color: var(--text-color);
  cursor: pointer;
  font-size: 16px;
  padding: 2px 4px;
  border-radius: 3px;
  opacity: 0.6;
  transition:
    opacity 0.2s,
    background-color 0.2s;
  margin-left: 4px;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.tab-close:hover {
  opacity: 1;
  background-color: var(--hover-bg);
}

.tab.active .tab-close {
  opacity: 0.8;
}
</style>
