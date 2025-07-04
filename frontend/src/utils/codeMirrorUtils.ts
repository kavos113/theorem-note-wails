import { EditorState } from '@codemirror/state';
import { EditorView } from '@codemirror/view';
import { basicSetup } from 'codemirror';
import { markdown, markdownLanguage } from '@codemirror/lang-markdown';
import { languages } from '@codemirror/language-data';
import { oneDark } from '@codemirror/theme-one-dark';
import { autocompletion, CompletionContext, CompletionResult } from '@codemirror/autocomplete';
import { LoadTheorems } from '../../wailsjs/go/main/App';

export interface CodeMirrorInstance {
  view: EditorView;
  updateContent: (content: string) => void;
  getContent: () => string;
  destroy: () => void;
  setEditorStyle: (style: Record<string, string>) => void;
}

const theoremAutocompletion = (rootDir: string) => {
  return autocompletion({
    override: [
      async (context: CompletionContext): Promise<CompletionResult | null> => {
        const match = context.matchBefore(/\[\[([^\]]*)$/);
        if (!match) {
          return null;
        }

        const theorems = await LoadTheorems(rootDir);
        const options = Object.entries(theorems).map(([name, path]) => ({
          label: name,
          apply: `${path}|${name}`,
          type: 'keyword'
        }));

        return {
          from: match.from + 2,
          options,
          validFor: /^\[\[[^\]]*$/
        };
      }
    ]
  });
};

export const createCodeMirrorEditor = (
  container: HTMLElement,
  initialContent: string,
  onChange: (content: string) => void,
  isDarkTheme = false,
  rootDir: string
): CodeMirrorInstance => {
  const editorTheme = EditorView.theme({
    '&': {
      height: '100%',
      fontSize: '14px',
      fontFamily: '"Consolas", "Monaco", "Courier New", monospace'
    },
    '.cm-content': {
      padding: '12px',
      minHeight: '100%'
    },
    '.cm-editor': {
      height: '100%'
    },
    '.cm-scroller': {
      fontFamily: 'inherit'
    },
    '.cm-line': {
      padding: '0 2px'
    },
    '.cm-gutters': {
      backgroundColor: 'var(--bg-color)',
      borderRight: '1px solid var(--border-color)'
    },
    '.cm-activeLine': {
      backgroundColor: 'var(--sidebar-bg)'
    },
    '.cm-activeLineGutter': {
      backgroundColor: 'var(--sidebar-bg)'
    },
    '.cm-selectionMatch': {
      backgroundColor: 'var(--accent-color-light)'
    },
    '.cm-link': {
      color: 'var(--accent-color)'
    },
    '.cm-url': {
      color: 'var(--accent-color)'
    },
    '.cm-keyword': {
      color: '#c678dd'
    },
    '.cm-comment': {
      color: '#5c6370',
      fontStyle: 'italic'
    },
    '.cm-string': {
      color: '#98c379'
    },
    '.cm-number': {
      color: '#d19a66'
    },
    '.cm-variableName': {
      color: '#e06c75'
    },
    '.cm-typeName': {
      color: '#e5c07b'
    },
    '.cm-operator': {
      color: '#56b6c2'
    },
    '.cm-punctuation': {
      color: '#abb2bf'
    },
    '.cm-property': {
      color: '#e06c75'
    },
    '.cm-className': {
      color: '#e5c07b'
    },
    '.cm-tag': {
      color: '#e06c75'
    },
    '.cm-attributeName': {
      color: '#d19a66'
    },
    '.cm-attributeValue': {
      color: '#98c379'
    },
    '.cm-qualifier': {
      color: '#e5c07b'
    },
    '.cm-meta': {
      color: '#abb2bf'
    },
    '.cm-hr': {
      color: '#abb2bf'
    },
    '.cm-quote': {
      color: '#98c379',
      fontStyle: 'italic'
    },
    '.cm-header': {
      color: '#61afef',
      fontWeight: 'bold'
    },
    '.cm-strong': {
      fontWeight: 'bold'
    },
    '.cm-emphasis': {
      fontStyle: 'italic'
    },
    '.cm-strikethrough': {
      textDecoration: 'line-through'
    },
    '.cm-atom': {
      color: '#d19a66'
    },
    '.cm-def': {
      color: '#e06c75'
    },
    '.cm-bracket': {
      color: '#abb2bf'
    },
    '.cm-builtin': {
      color: '#e5c07b'
    },
    '.cm-error': {
      color: '#e06c75'
    },
    '.cm-invalid': {
      color: '#e06c75'
    },
    '.cm-codeblock': {
      backgroundColor: '#2c313a',
      borderRadius: '4px',
      padding: '10px',
      fontFamily: '"Fira Code", "Courier New", monospace'
    }
  });

  const extensions = [
    basicSetup,
    markdown({
      base: markdownLanguage,
      codeLanguages: languages
    }),
    EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        onChange(update.state.doc.toString());
      }
    }),
    editorTheme,
    theoremAutocompletion(rootDir)
  ];

  if (isDarkTheme) {
    extensions.push(oneDark);
  }

  const state = EditorState.create({
    doc: initialContent,
    extensions
  });

  const view = new EditorView({
    state,
    parent: container
  });

  return {
    view,
    updateContent: (content: string) => {
      const currentContent = view.state.doc.toString();
      if (currentContent !== content) {
        view.dispatch({
          changes: {
            from: 0,
            to: view.state.doc.length,
            insert: content
          }
        });
      }
    },
    getContent: () => view.state.doc.toString(),
    destroy: () => view.destroy(),
    setEditorStyle: (style: Record<string, string>) => {
      // CodeMirrorのスタイルを動的に更新
      // EditorView.themeは一度作成すると変更できないため、
      // 直接DOMを操作するか、新しいテーマを適用する必要があります。
      // ここでは簡易的にDOMを操作します。
      const cmElement = view.dom;
      if (cmElement) {
        if (style.fontFamily) {
          cmElement.style.fontFamily = style.fontFamily;
        }
        if (style.fontSize) {
          cmElement.style.fontSize = style.fontSize;
        }
      }
    }
  };
};
