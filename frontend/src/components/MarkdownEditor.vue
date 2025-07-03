<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick, computed } from 'vue';
import { createCodeMirrorEditor, type CodeMirrorInstance } from '../utils/codeMirrorUtils';
import type { ViewMode } from '../types/viewMode';
import '../assets/styles/highlight.css';
import '../assets//styles/katex.css';
import { markdownToHtml, getProjectRoot, renderMermaid } from '../utils/markdownUtils';
import 'highlight.js/styles/github.css';
import { WriteFile, GetFontSettings } from '../../wailsjs/go/main/App';
import { EventsOn } from '../../wailsjs/runtime';
import type { backend } from '../../wailsjs/go/models';

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

const handleInternalLinkClick = (event: MouseEvent) => {
  const target = event.target as HTMLElement;
  const link = target.closest('a[data-internal-link="true"]');

  if (link) {
    event.preventDefault();
    const path = link.getAttribute('data-path');
    const header = link.getAttribute('data-header');

    if (path) {
      // @ts-expect-error header is not defined in the event type
      emit('select-file', path, header);
    }
  }
};

const localContent = ref(props.fileContent);
const isSaving = ref(false);
const editorWidth = ref(50); // エディタの幅(%)
const isResizing = ref(false);
const editorContainer = ref<HTMLElement>();
const codeMirrorInstance = ref<CodeMirrorInstance>();
const previewContainer = ref<HTMLElement>();

const htmlPreview = ref<string>('');

// スクロール同期のためのフラグ
const isSyncingScroll = ref(false);
const fontSettings = ref<backend.FontSettings>();

// 動的なスタイルを計算
const editorStyle = computed(() => ({
  fontFamily: fontSettings.value?.editor_font_family || 'monospace',
  fontSize: `${fontSettings.value?.editor_font_size || 14}px`
}));

const previewStyle = computed(() => ({
  fontFamily: fontSettings.value?.preview_font_family || 'sans-serif',
  fontSize: `${fontSettings.value?.preview_font_size || 14}px`
}));

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
  console.log('Content changed:', event);
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
    if (newContent === null || newContent === undefined) {
      htmlPreview.value = '';
      return;
    }
    try {
      htmlPreview.value = await markdownToHtml(newContent);
      await nextTick();
      renderMermaid();
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
  setupScrollListeners();

  // フォント設定を読み込み、変更を監視
  await loadFontSettings();
  cleanupFontListener = EventsOn('font-settings-updated', (settings: backend.FontSettings) => {
    applyFontSettings(settings);
  });

  setupLinkListener();
});

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeyDown);

  if (codeMirrorInstance.value) {
    codeMirrorInstance.value.destroy();
  }
  // スクロールイベントリスナーを削除
  removeScrollListeners();
  // イベントリスナーをクリーンアップ
  if (cleanupFontListener) {
    cleanupFontListener();
  }
  removeLinkListener();
});

// --- フォント設定 ---
let cleanupFontListener: () => void;

const applyFontSettings = (settings: backend.FontSettings) => {
  fontSettings.value = settings;
  if (codeMirrorInstance.value) {
    codeMirrorInstance.value.setEditorStyle(editorStyle.value);
  }
};


const setupLinkListener = () => {
  if (previewContainer.value) {
    previewContainer.value.addEventListener('click', handleInternalLinkClick);
  }
};

const removeLinkListener = () => {
  if (previewContainer.value) {
    previewContainer.value.removeEventListener('click', handleInternalLinkClick);
  }
};

const loadFontSettings = async () => {
  const rootDir = getProjectRoot();
  if (!rootDir) return;
  try {
    const settings = await GetFontSettings(rootDir);
    applyFontSettings(settings);
  } catch (err) {
    console.error('フォント設定の読み込みに失敗しました:', err);
  }
};

// スクロール同期ロジック
let editorScroller: HTMLElement | null = null;

const handleEditorScroll = () => {
  if (isSyncingScroll.value || !editorScroller || !previewContainer.value) return;
  isSyncingScroll.value = true;

  const editorScrollRatio =
    editorScroller.scrollTop / (editorScroller.scrollHeight - editorScroller.clientHeight);
  previewContainer.value.scrollTop =
    editorScrollRatio * (previewContainer.value.scrollHeight - previewContainer.value.clientHeight);

  requestAnimationFrame(() => {
    isSyncingScroll.value = false;
  });
};

const handlePreviewScroll = () => {
  if (isSyncingScroll.value || !editorScroller || !previewContainer.value) return;
  isSyncingScroll.value = true;

  const previewScrollRatio =
    previewContainer.value.scrollTop /
    (previewContainer.value.scrollHeight - previewContainer.value.clientHeight);
  editorScroller.scrollTop =
    previewScrollRatio * (editorScroller.scrollHeight - editorScroller.clientHeight);

  requestAnimationFrame(() => {
    isSyncingScroll.value = false;
  });
};

const setupScrollListeners = () => {
  if (props.viewMode !== 'split') return;

  nextTick(() => {
    if (codeMirrorInstance.value && previewContainer.value) {
      editorScroller = editorContainer.value?.querySelector('.cm-scroller') as HTMLElement;
      if (editorScroller) {
        editorScroller.addEventListener('scroll', handleEditorScroll);
        previewContainer.value.addEventListener('scroll', handlePreviewScroll);
      }
    }
  });
};

const removeScrollListeners = () => {
  if (editorScroller) {
    editorScroller.removeEventListener('scroll', handleEditorScroll);
  }
  if (previewContainer.value) {
    previewContainer.value.removeEventListener('scroll', handlePreviewScroll);
  }
};

watch(
  () => props.viewMode,
  (newMode) => {
    removeScrollListeners();
    if (newMode === 'split') {
      setupScrollListeners();
    }
  }
);

watch(
  () => props.selectedFilePath,
  () => {
    removeScrollListeners();
    if (props.viewMode === 'split') {
      setupScrollListeners();
    }
  }
);

const scrollToHeader = (header: string) => {
  if (!previewContainer.value) return;

  const headerId = header
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '-');

  const headerElement = previewContainer.value.querySelector(`#${headerId}`);

  if (headerElement) {
    headerElement.scrollIntoView({ behavior: 'smooth' });
  } else {
    console.warn(`Header with id #${headerId} not found.`);
  }
};

defineExpose({
  scrollToHeader
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
        <div ref="editorContainer" class="codemirror-container" :style="editorStyle"></div>
        <textarea
          v-model="localContent"
          class="markdown-editor hidden"
          :style="editorStyle"
          @input="handleContentChange"
        ></textarea>
      </div>
      <div v-if="viewMode === 'split'" class="resizer" @mousedown="startResize"></div>
      <div v-if="viewMode !== 'editor'" class="preview-pane" :style="{ width: getPreviewWidth() }">
        <div class="pane-header">プレビュー</div>
        <!-- eslint-disable-next-line vue/no-v-html -->
        <div
          ref="previewContainer"
          class="markdown-preview"
          :style="previewStyle"
          v-html="htmlPreview"
        ></div>
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

.markdown-preview :deep(.mermaid) {
  background-color: var(--bg-color);
  padding: 1em;
  border-radius: 5px;
  margin: 1em 0;
  text-align: center;
}

.markdown-preview :deep(.card-link-container) {
  width: 100%;
}

.markdown-preview :deep(.card-link) {
  display: flex;
  background-color: #ffffff;
  border: 1px solid #b1b8bd;
  border-radius: 8px;
  overflow: hidden;
  text-decoration: none;
  color: #14171a;
  transition:
    transform 0.2s ease-in-out,
    box-shadow 0.2s ease-in-out;
}

.markdown-preview :deep(.card-link:hover) {
  background-color: #f0f0f0;
}

.markdown-preview :deep(.card-content) {
  flex: 1;
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0; /* Flexboxでのテキストオーバーフロー問題を防止 */
}

.markdown-preview :deep(.card-title) {
  font-weight: 600;
  margin: 0 0 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.markdown-preview :deep(.card-description) {
  color: #657786;
  margin: 0 0 12px;
  flex-grow: 1;
  /* 説明文を2行に制限して、...で省略 */
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
}

.markdown-preview :deep(.card-footer) {
  display: flex;
  align-items: center;
  color: #657786;
}

.markdown-preview :deep(.card-favicon) {
  width: 16px;
  height: 16px;
  margin-right: 8px;
  flex-shrink: 0;
}

.markdown-preview :deep(.card-url) {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.markdown-preview :deep(.card-thumbnail) {
  width: 130px;
  flex-shrink: 0;
  background-color: #f5f8fa;
}

.markdown-preview :deep(.card-thumbnail) img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-left: 1px solid #e1e8ed;
}
</style>
