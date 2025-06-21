import { EditorState } from '@codemirror/state';
import { EditorView } from '@codemirror/view';
import { basicSetup } from 'codemirror';
import { markdown } from '@codemirror/lang-markdown';
import { oneDark } from '@codemirror/theme-one-dark';

export interface CodeMirrorInstance {
  view: EditorView;
  updateContent: (content: string) => void;
  getContent: () => string;
  destroy: () => void;
}

export const createCodeMirrorEditor = (
  container: HTMLElement,
  initialContent: string,
  onChange: (content: string) => void,
  isDarkTheme = false
): CodeMirrorInstance => {
  const extensions = [
    basicSetup,
    markdown(),
    EditorView.updateListener.of((update) => {
      if (update.docChanged) {
        onChange(update.state.doc.toString());
      }
    }),
    EditorView.theme({
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
      }
    })
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
    destroy: () => view.destroy()
  };
};
