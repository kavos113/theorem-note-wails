import { unified } from 'unified';
import remarkParse from 'remark-parse';
import remarkRehype from 'remark-rehype';
import rehypeStringify from 'rehype-stringify';
import rehypeHighlight from 'rehype-highlight';
import remarkBreaks from 'remark-breaks';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeKatex from 'rehype-katex';
import { visit } from 'unist-util-visit';
import mermaid from 'mermaid';
import type { Root, Element } from 'hast';

let projectRoot = '';

const IMAGE_PREFIX = '/_images/';

export const setProjectRoot = (root: string): void => {
  if (root) {
    projectRoot = root.replace(/\\/g, '/');
  } else {
    console.warn('Project root is not set. Using default empty string.');
  }
};

export const getProjectRoot = (): string => {
  return projectRoot;
};

const rehypeMermaid = () => {
  return (tree: Root) => {
    visit(tree, 'element', (node: Element, index?: number, parent?: Root | Element) => {
      if (node.tagName === 'pre' && parent?.children && index !== undefined) {
        const codeNode = node.children[0];
        if (
          codeNode?.type === 'element' &&
          codeNode.tagName === 'code' &&
          (codeNode.properties?.className as string[])?.includes('language-mermaid')
        ) {
          const textNode = codeNode.children[0];
          if (textNode?.type === 'text') {
            const code = textNode.value;
            const mermaidContainer: Element = {
              type: 'element',
              tagName: 'div',
              properties: { className: ['mermaid'] },
              children: [{ type: 'text', value: code }]
            };
            parent.children.splice(index, 1, mermaidContainer);
          }
        }
      }
    });
  };
};

export const markdownToHtml = async (markdown: string): Promise<string> => {
  const parsed = await unified()
    .use(remarkParse)
    .use(remarkBreaks)
    .use(remarkGfm)
    .use(remarkMath)
    .use(remarkRehype, { allowDangerousHtml: true })
    .use(rehypeKatex)
    .use(rehypeHighlight)
    .use(rehypeMermaid)
    .use(rehypeStringify)
    .process(convertObsidianLinks(markdown));

  return parsed.toString();
};

function convertObsidianLinks(markdown: string): string {
  console.log(projectRoot);
  const regex = /!\[\[(.*?)]]/g;
  return markdown.replace(regex, (_, filename) => {
    const encoded = encodeURIComponent(filename);
    return `![${filename}](${projectRoot}${IMAGE_PREFIX}${encoded})`;
  });
}

export const renderMermaid = () => {
  mermaid.initialize({ startOnLoad: false });
  mermaid.run({
    nodes: document.querySelectorAll('.mermaid')
  });
};
