<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { backend } from '../../wailsjs/go/models';
import FileItem = backend.FileItem;

const props = defineProps<{
  item: FileItem;
  selectedItem?: string | null;
  expandedItems?: Set<string>;
}>();

const emit = defineEmits<{
  (e: 'select-file', path: string, header?: string): void;
  (e: 'select-item', path: string): void;
  (e: 'expand-item', path: string): void;
}>();

const isExpanded = ref(false);

const isSelected = computed(() => {
  return props.selectedItem === props.item.Path;
});

const handleClick = (): void => {
  emit('select-item', props.item.Path);
  if (props.item.IsDirectory) {
    isExpanded.value = !isExpanded.value;
    emit('expand-item', props.item.Path);
  } else {
    emit('select-file', props.item.Path, undefined);
  }
};

onMounted(() => {
  if (props.expandedItems && props.expandedItems.has(props.item.Path)) {
    isExpanded.value = true;
  }
});
</script>

<template>
  <li class="file-tree-item">
    <div
      class="item-content"
      :class="{ 'is-directory': item.IsDirectory, 'is-selected': isSelected }"
      @click="handleClick"
    >
      <span v-if="item.IsDirectory" class="folder-arrow" :class="{ expanded: isExpanded }">
        &gt;
      </span>
      <span v-else class="folder-arrow"></span>
      <span class="icon">
        <span v-if="item.IsDirectory" class="folder-icon">
          {{ isExpanded ? '📂' : '📁' }}
        </span>
        <span v-else class="file-icon">📄</span>
      </span>
      <span class="name">{{ item.Name }}</span>
    </div>
    <ul v-if="item.IsDirectory && isExpanded && item.Children" class="children">
      <file-tree-item
        v-for="child in item.Children"
        :key="child.Path"
        :item="child"
        :selected-item="selectedItem"
        :expanded-items="expandedItems"
        @select-file="$emit('select-file', $event)"
        @select-item="$emit('select-item', $event)"
        @expand-item="$emit('expand-item', $event)"
      />
    </ul>
  </li>
</template>

<style scoped>
.file-tree-item {
  list-style-type: none;
  margin: 0;
  padding: 0;
}

.item-content {
  display: flex;
  align-items: center;
  padding: 0 8px;
  cursor: pointer;
  border-radius: 3px;
}

.item-content:hover {
  background-color: var(--hover-bg, #eaeaea);
}

.is-selected {
  background-color: var(--accent-color, #0078d7);
  color: white;
}

.is-selected:hover {
  background-color: var(--accent-color, #0078d7);
}

.folder-arrow {
  display: inline-block;
  width: 12px;
  height: 12px;
  margin-right: 4px;
  font-size: 10px;
  transition: transform 0.2s;
  transform: rotate(0deg);
}

.folder-arrow.expanded {
  transform: rotate(90deg);
}

.icon {
  margin-right: 6px;
  font-size: 14px;
  width: 16px;
  text-align: center;
}

.name {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.children {
  list-style-type: none;
  margin: 0;
  padding: 0 0 0 16px;
}
</style>
