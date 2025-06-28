<script setup lang="ts">
import { ref, watch } from 'vue';
import { GetFontSettings, SaveFontSettings } from '../../wailsjs/go/main/App';
import type { backend } from '../../wailsjs/go/models';

interface Props {
  isOpen: boolean;
  rootDir: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{ (e: 'close'): void }>();

const fontSettings = ref<backend.FontSettings>({} as backend.FontSettings);

const loadSettings = async () => {
  if (!props.rootDir) return;
  try {
    const settings = await GetFontSettings(props.rootDir);
    fontSettings.value = settings;
  } catch (err) {
    console.error('フォント設定の読み込みに失敗しました:', err);
  }
};

const saveSettings = async () => {
  if (!props.rootDir) return;
  try {
    await SaveFontSettings(props.rootDir, fontSettings.value);
    emit('close');
  } catch (err) {
    console.error('フォント設定の保存に失敗しました:', err);
  }
};

watch(
  () => props.isOpen,
  (newVal) => {
    if (newVal) {
      loadSettings();
    }
  }
);
</script>

<template>
  <div v-if="isOpen" class="modal-overlay" @click.self="emit('close')">
    <div class="modal-content">
      <h2>設定</h2>

      <div class="settings-group">
        <h3>エディタのフォント</h3>
        <div class="form-item">
          <label for="editor-font-family">フォントファミリー</label>
          <input id="editor-font-family" v-model="fontSettings.editor_font_family" type="text" />
        </div>
        <div class="form-item">
          <label for="editor-font-size">フォントサイズ (px)</label>
          <input
            id="editor-font-size"
            v-model.number="fontSettings.editor_font_size"
            type="number"
          />
        </div>
      </div>

      <div class="settings-group">
        <h3>プレビューのフォント</h3>
        <div class="form-item">
          <label for="preview-font-family">フォントファミリー</label>
          <input id="preview-font-family" v-model="fontSettings.preview_font_family" type="text" />
        </div>
        <div class="form-item">
          <label for="preview-font-size">フォントサイズ (px)</label>
          <input
            id="preview-font-size"
            v-model.number="fontSettings.preview_font_size"
            type="number"
          />
        </div>
      </div>

      <div class="modal-actions">
        <button class="btn-secondary" @click="emit('close')">キャンセル</button>
        <button class="btn-primary" @click="saveSettings">保存</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: var(--bg-color);
  padding: 20px;
  border-radius: 8px;
  width: 500px;
  max-width: 90%;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

h2 {
  margin-top: 0;
  margin-bottom: 20px;
  border-bottom: 1px solid var(--border-color);
  padding-bottom: 10px;
}

.settings-group {
  margin-bottom: 20px;
}

.settings-group h3 {
  margin-bottom: 15px;
  font-size: 1.1em;
}

.form-item {
  margin-bottom: 10px;
}

.form-item label {
  display: block;
  margin-bottom: 5px;
  font-size: 14px;
}

.form-item input {
  width: 100%;
  padding: 8px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background-color: var(--bg-color);
  color: var(--text-color);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.btn-primary,
.btn-secondary {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.btn-primary {
  background-color: var(--accent-color);
  color: white;
}

.btn-secondary {
  background-color: var(--sidebar-header-bg);
  color: var(--text-color);
}
</style>
