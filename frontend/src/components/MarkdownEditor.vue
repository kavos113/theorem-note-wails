<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue';
import { createCodeMirrorEditor, type CodeMirrorInstance } from '../utils/codeMirrorUtils';
import type { ViewMode } from '../types/viewMode';
import '../assets/styles/highlight.css';
import '../assets//styles/katex.css';
import { markdownToHtml } from '../utils/markdownUtils';
import 'highlight.js/styles/github.css';
import { WriteFile } from "../../wailsjs/go/main/App";

interface Props {
  selectedFilePath?: string;
  fileContent: string;
  viewMode: ViewMode;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'update:fileContent', value: string): void;
  (e: 'file-saved'): void;
}

const emit = defineEmits<Emits>();

const localContent = ref(props.fileContent);
const isSaving = ref(false);
const editorWidth = ref(50); // エディタの幅(%)
const isResizing = ref(false);
const editorContainer = ref<HTMLElement>();
const codeMirrorInstance = ref<CodeMirrorInstance>();

const htmlPreview = ref<string>('');

const saveFile = async (): Promise<void> => {
  if (!props.selectedFilePath || isSaving.value) return;

  try {
    isSaving.value = true;
    await WriteFile(props.selectedFilePath, localContent.value);
    emit('file-saved');
    console.log('ファイルが保存されました');
  } catch (err) {
    console.error('ファイル保存エラー:', err);
  } finally {
    isSaving.value = false;
  }
};

watch(
  () => props.fileContent,
  (newContent) => {
    localContent.value = newContent;
    // CodeMirrorの内容も更新
    if (codeMirrorInstance.value) {
      codeMirrorInstance.value.updateContent(newContent);
    }
  },
  { immediate: true }
);

const handleContentChange = (event: Event): void => {
  const target = event.target as HTMLTextAreaElement;
  localContent.value = target.value;
  emit('update:fileContent', target.value);
};

const handleCodeMirrorChange = (content: string): void => {
  localContent.value = content;
  emit('update:fileContent', content);
};

watch(
  localContent,
  async (newContent) => {
    if (!newContent) return;
    try {
      htmlPreview.value = await markdownToHtml(newContent);
    } catch (err) {
      console.error('マークダウン変換エラー:', err);
      htmlPreview.value = '<p>プレビューの生成中にエラーが発生しました</p>';
    }
  },
  { immediate: true }
);

const handleKeyDown = (event: KeyboardEvent): void => {
  if (event.ctrlKey && event.key === 's') {
    event.preventDefault();
    saveFile();
  }
};

const startResize = (event: MouseEvent): void => {
  event.preventDefault();
  isResizing.value = true;

  const target = event.currentTarget as HTMLElement;
  const container = target.parentElement as HTMLElement;
  const containerRect = container.getBoundingClientRect();

  const handleMouseMove = (e: MouseEvent): void => {
    if (!isResizing.value) return;

    const newWidth = ((e.clientX - containerRect.left) / containerRect.width) * 100;
    editorWidth.value = Math.min(Math.max(newWidth, 20), 80); // 20%-80%の範囲で制限
  };

  const handleMouseUp = (): void => {
    isResizing.value = false;
    document.removeEventListener('mousemove', handleMouseMove);
    document.removeEventListener('mouseup', handleMouseUp);
  };

  document.addEventListener('mousemove', handleMouseMove);
  document.addEventListener('mouseup', handleMouseUp);
};

const getEditorWidth = (): string => {
  if (props.viewMode === 'editor') return '100%';
  if (props.viewMode === 'split') return `${editorWidth.value}%`;
  return '0%';
};

const getPreviewWidth = (): string => {
  if (props.viewMode === 'preview') return '100%';
  if (props.viewMode === 'split') return `${100 - editorWidth.value}%`;
  return '0%';
};

onMounted(async () => {
  document.addEventListener('keydown', handleKeyDown);

  await nextTick();
  if (editorContainer.value) {
    codeMirrorInstance.value = createCodeMirrorEditor(
      editorContainer.value,
      localContent.value,
      handleCodeMirrorChange
    );
  }
});

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown);

  if (codeMirrorInstance.value) {
    codeMirrorInstance.value.destroy();
  }
});
</script>

<template>
  <div class="markdown-editor-container">
    <div class="editor-header">
      {{ selectedFilePath }}
    </div>
    <div class="editor-content-split" :class="`view-mode-${viewMode}`">
      <div v-if="viewMode !== 'preview'" class="editor-pane" :style="{ width: getEditorWidth() }">
        <div class="pane-header">エディタ</div>
        <div ref="editorContainer" class="codemirror-container"></div>
        <textarea
          v-model="localContent"
          class="markdown-editor hidden"
          @input="handleContentChange"
        ></textarea>
      </div>
      <div v-if="viewMode === 'split'" class="resizer" @mousedown="startResize"></div>
      <div v-if="viewMode !== 'editor'" class="preview-pane" :style="{ width: getPreviewWidth() }">
        <div class="pane-header">プレビュー</div>
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div class="markdown-preview" v-html="htmlPreview"></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.markdown-editor-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.editor-header {
  padding: 5px 10px;
  font-size: 12px;
  color: var(--text-color);
  background-color: var(--bg-color);
  border-bottom: 1px solid var(--border-color);
}

.editor-content-split {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.editor-pane {
  display: flex;
  flex-direction: column;
  min-width: 200px;
}

.resizer {
  width: 4px;
  background-color: var(--border-color);
  cursor: col-resize;
  flex-shrink: 0;
  transition: background-color 0.2s ease;
}

.resizer:hover {
  background-color: var(--accent-color);
}

.preview-pane {
  display: flex;
  flex-direction: column;
  min-width: 200px;
}

.pane-header {
  padding: 8px 12px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-color);
  background-color: var(--sidebar-header-bg);
  border-bottom: 1px solid var(--border-color);
}

.markdown-editor {
  flex: 1;
  padding: 12px;
  border: none;
  outline: none;
  background-color: var(--bg-color);
  color: var(--text-color);
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  resize: none;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.codemirror-container {
  flex: 1;
  height: 100%;
  overflow: hidden;
}

.hidden {
  display: none;
}

.markdown-preview {
  flex: 1;
  padding: 12px;
  overflow-y: auto;
  background-color: var(--bg-color);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-color);
}

.editor-content-split.preview-hidden .editor-pane {
  width: 100% !important;
}

.view-mode-editor .editor-pane {
  width: 100% !important;
}

.view-mode-preview .preview-pane {
  width: 100% !important;
}

/* マークダウンプレビューのスタイリング */
.markdown-preview :deep(h1) {
  font-size: 2em;
  font-weight: bold;
  margin: 0.67em 0;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.3em;
}

.markdown-preview :deep(h2) {
  font-size: 1.5em;
  font-weight: bold;
  margin: 0.83em 0;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 0.3em;
}

.markdown-preview :deep(h3) {
  font-size: 1.17em;
  font-weight: bold;
  margin: 1em 0;
}

.markdown-preview :deep(h4) {
  font-size: 1em;
  font-weight: bold;
  margin: 1.33em 0;
}

.markdown-preview :deep(h5) {
  font-size: 0.83em;
  margin: 1.67em 0;
}

.markdown-preview :deep(h6) {
  font-size: 0.67em;
  margin: 2.33em 0;
}

.markdown-preview :deep(p) {
  margin: 1em 0;
}

.markdown-preview :deep(blockquote) {
  margin: 1em 0;
  padding-left: 1em;
  border-left: 4px solid var(--accent-color);
  background-color: var(--sidebar-bg);
  color: #666;
}

.markdown-preview :deep(ul),
.markdown-preview :deep(ol) {
  margin: 1em 0;
  padding-left: 2em;
}

.markdown-preview :deep(ul) {
  list-style-type: disc;
}

.markdown-preview :deep(ol) {
  list-style-type: decimal;
}

.markdown-preview :deep(code) {
  background-color: var(--sidebar-bg);
  padding: 0.2em 0.4em;
  border-radius: 3px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.9em;
}

.markdown-preview :deep(pre) {
  background-color: var(--sidebar-bg);
  padding: 1em;
  border-radius: 5px;
  overflow-x: auto;
  margin: 1em 0;
}

.markdown-preview :deep(pre code) {
  background-color: transparent;
  padding: 0;
  border-radius: 0;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.9em;
}

.markdown-preview :deep(table) {
  border-collapse: collapse;
  margin: 1em 0;
}

.markdown-preview :deep(th),
.markdown-preview :deep(td) {
  border: 1px solid var(--border-color);
  padding: 0.5em;
  text-align: left;
}

.markdown-preview :deep(th) {
  background-color: var(--sidebar-header-bg);
  font-weight: bold;
}

.markdown-preview :deep(strong) {
  font-weight: bold;
}
</style>
