<template>
  <div class="markdown-content">
    <div v-html="markdownToHtml(text)"></div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import { marked } from 'marked';
import DOMPurify from 'dompurify';

export default defineComponent({
  name: 'MarkDownComponent',
  props: {
    text: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    // Парсим Markdown в HTML и санитизируем
    const markdownToHtml = (markdown: string): string => {
      const rawHtml = marked.parse(markdown, {
        gfm: true, // Поддержка GitHub Flavored Markdown
        breaks: true, // Преобразование \n в <br>
      }) as string;
      return DOMPurify.sanitize(rawHtml);
    };

    return {
      markdownToHtml,
    };
  },
});
</script>

<style scoped>
/* Базовые стили для prose */
.markdown-content {
  @apply max-w-3xl mx-auto p-1;
}

/* Стили для заголовков */
.markdown-content :deep(h1) {
  @apply text-2xl font-bold text-gray-900 mb-6;
}
.markdown-content :deep(h2) {
  @apply text-xl font-semibold text-gray-800 mb-4;
}

/* Стили для текста */
.markdown-content :deep(p) {
  @apply text-base text-gray-700 leading-7 mb-4;
}
.markdown-content :deep(li) {
  @apply text-base text-gray-700 leading-7 mb-2 list-disc list-outside ml-6;
}
.markdown-content :deep(code) {
  @apply bg-gray-100 px-1 py-0.5 rounded text-sm text-gray-800;
}

/* Стили для списков */
.markdown-content :deep(ul) {
  @apply list-disc list-outside ml-6 mb-4;
}
.markdown-content :deep(ol) {
  @apply list-decimal list-outside ml-6 mb-4;
}
</style>