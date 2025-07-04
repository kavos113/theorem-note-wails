import { unified } from 'unified';
import remarkParse from 'remark-parse';
import remarkRehype from 'remark-rehype';
import rehypeStringify from 'rehype-stringify';
import rehypeHighlight from 'rehype-highlight';
import remarkBreaks from 'remark-breaks';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeKatex from 'rehype-katex';
import rehypeSlug from 'rehype-slug';
import { visit } from 'unist-util-visit';
import mermaid from 'mermaid';
import type { Element, Root } from 'hast';
import type { Root as RemarkRoot, Text } from 'mdast';
import rehypeRaw from 'rehype-raw';

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

interface CardLinkData {
  url: string;
  title: string;
  description?: string;
  host: string;
  favicon?: string;
  image?: string;
}

const makeCardLinkElement = (data: CardLinkData): Element => {
  const children: Element[] = [
    {
      type: 'element',
      tagName: 'div',
      properties: { className: ['card-content'] },
      children: [
        {
          type: 'element',
          tagName: 'p',
          properties: { className: ['card-title'] },
          children: [{ type: 'text', value: data.title }]
        },
        {
          type: 'element',
          tagName: 'p',
          properties: { className: ['card-description'] },
          children: [{ type: 'text', value: data.description || '' }]
        },
        {
          type: 'element',
          tagName: 'div',
          properties: { className: ['card-footer'] },
          children: [
            {
              type: 'element',
              tagName: 'img',
              properties: {
                src: data.favicon,
                alt: 'site favicon',
                className: ['card-favicon']
              },
              children: []
            },
            {
              type: 'element',
              tagName: 'span',
              properties: { className: ['card-url'] },
              children: [{ type: 'text', value: data.host }]
            }
          ]
        }
      ]
    }
  ];

  if (data.image) {
    children.push({
      type: 'element',
      tagName: 'div',
      properties: { className: ['card-thumbnail'] },
      children: [
        {
          type: 'element',
          tagName: 'img',
          properties: {
            src: data.image,
            alt: data.title + ' Logo'
          },
          children: []
        }
      ]
    });
  }

  return {
    type: 'element',
    tagName: 'div',
    properties: { className: ['card-link-container'] },
    children: [
      {
        type: 'element',
        tagName: 'a',
        properties: {
          href: data.url,
          className: ['card-link'],
          target: '_blank',
          rel: 'noopener noreferrer'
        },
        children: children
      }
    ]
  };
};

const rehypeCardLink = () => {
  return (tree: Root) => {
    visit(tree, 'element', (node: Element, index?: number, parent?: Root | Element) => {
      if (node.tagName === 'pre' && parent?.children && index !== undefined) {
        const codeNode = node.children[0];
        if (
          codeNode?.type === 'element' &&
          codeNode.tagName === 'code' &&
          (codeNode.properties?.className as string[])?.includes('language-cardlink')
        ) {
          const textNode = codeNode.children[0];
          if (textNode?.type === 'text') {
            const code = textNode.value;
            const lines = code.split('\n');
            const data: Partial<CardLinkData> = {};
            for (const line of lines) {
              const [key, ...rest] = line.split(':');
              if (!key || rest.length === 0) continue;
              const trimmedKey = key.trim() as keyof CardLinkData;
              data[trimmedKey] = rest.join(':').trim().replace(/^"|"$/g, '');
            }

            parent.children.splice(index, 1, makeCardLinkElement(data as CardLinkData));
          }
        }
      }
    });
  };
};

const remarkInternalLinks = () => {
  return (tree: RemarkRoot) => {
    visit(tree, 'text', (node: Text, index, parent) => {
      const regex = /\[\[(.*?)]]/g;
      let match;
      let lastIndex = 0;
      const newNodes: any[] = [];
      const value = node.value;

      while ((match = regex.exec(value)) !== null) {
        if (match.index > lastIndex) {
          newNodes.push({
            type: 'text',
            value: value.slice(lastIndex, match.index)
          });
        }
        const [path, header] = match[1].split('|');
        const text = header ? `${path}#${header}` : path;
        newNodes.push({
          type: 'link',
          url: '#',
          data: {
            hProperties: {
              'data-internal-link': 'true',
              'data-path': path,
              ...(header ? { 'data-header': header } : {})
            }
          },
          children: [{ type: 'text', value: text }]
        });
        lastIndex = match.index + match[0].length;
      }
      if (lastIndex < value.length) {
        newNodes.push({
          type: 'text',
          value: value.slice(lastIndex)
        });
      }
      if (newNodes.length > 0 && parent && typeof index === 'number') {
        parent.children.splice(index, 1, ...newNodes);
        return index + newNodes.length;
      }
    });
  };
};

const rehypeTheoremMarkdown = () => {
  return (tree: Root) => {
    visit(tree, 'element', (node: Element, index?: number, parent?: Root | Element) => {
      if (node.tagName === 'theorem' && parent?.children && index !== undefined) {
        // theorem要素のname属性を取得
        const nameAttr = node.properties?.name as string;

        // theorem内のテキストコンテンツを収集
        let content = '';
        const collectText = (children: any[]) => {
          for (const child of children) {
            if (child.type === 'text') {
              content += child.value;
              console.log(child.value);
            } else if (child.children) {
              collectText(child.children);
            }
          }
        };
        collectText(node.children);

        const processor = unified()
          .use(remarkParse)
          .use(remarkBreaks)
          .use(remarkGfm)
          .use(remarkMath)
          .use(remarkRehype, { allowDangerousHtml: true });

        const mdast = processor.parse(content);
        const hast = processor.runSync(mdast);

        console.log(hast);

        parent.children[index] = {
          type: 'element',
          tagName: 'div',
          properties: {
            className: ['theorem']
          },
          children: [
            {
              type: 'element',
              tagName: 'h4',
              properties: { className: ['theorem-title'] },
              children: [{ type: 'text', value: nameAttr || 'Theorem' }]
            },
            {
              type: 'element',
              tagName: 'div',
              properties: { className: ['theorem-content'] },
              children: hast.children as Element[]
            }
          ]
        };
      }
    });
  };
};

export const markdownToHtml = async (markdown: string): Promise<string> => {
  const parsed = await unified()
    .use(remarkParse)
    .use(remarkBreaks)
    .use(remarkGfm)
    .use(remarkInternalLinks)
    .use(remarkMath)
    .use(remarkRehype, { allowDangerousHtml: true })
    .use(rehypeRaw)
    .use(rehypeTheoremMarkdown)
    .use(rehypeSlug)
    .use(rehypeKatex)
    .use(rehypeHighlight)
    .use(rehypeMermaid)
    .use(rehypeCardLink)
    .use(rehypeStringify)
    .process(convertObsidianLinks(markdown));

  return parsed.toString();
};

function convertObsidianLinks(markdown: string): string {
  const regex = /!\[\[(.*?)]]/g;
  return markdown.replace(regex, (_, filename) => {
    const encoded = encodeURIComponent(filename);
    return `![${filename}](file:///${projectRoot}${IMAGE_PREFIX}${encoded})`;
  });
}

export const renderMermaid = () => {
  mermaid.initialize({ startOnLoad: false });
  mermaid.run({
    nodes: document.querySelectorAll('.mermaid')
  });
};
