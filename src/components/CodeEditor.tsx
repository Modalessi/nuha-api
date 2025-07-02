import { Editor } from '@monaco-editor/react'
import { useState } from 'react'

interface CodeEditorProps {
  value: string
  onChange: (value: string) => void
  language: string
  height?: string
}

const languageMap: Record<string, string> = {
  '71': 'python',
  '92': 'python', 
  '100': 'python',
  '62': 'java',
  '54': 'cpp',
  '53': 'cpp',
  '52': 'cpp',
  '50': 'c',
  '49': 'c',
  '48': 'c',
  '63': 'javascript',
  '74': 'typescript',
}

export default function CodeEditor({ value, onChange, language, height = '400px' }: CodeEditorProps) {
  const [theme, setTheme] = useState<'light' | 'dark'>('light')
  
  const editorLanguage = languageMap[language] || 'python'

  const handleEditorChange = (value: string | undefined) => {
    onChange(value || '')
  }

  return (
    <div className="border border-gray-300 rounded-lg overflow-hidden">
      <div className="bg-gray-50 px-4 py-2 border-b border-gray-200 flex justify-between items-center">
        <span className="text-sm font-medium text-gray-700">
          Code Editor ({editorLanguage})
        </span>
        <button
          onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
          className="text-xs px-2 py-1 bg-gray-200 hover:bg-gray-300 rounded transition-colors"
        >
          {theme === 'light' ? '🌙' : '☀️'}
        </button>
      </div>
      <Editor
        height={height}
        language={editorLanguage}
        value={value}
        onChange={handleEditorChange}
        theme={theme === 'dark' ? 'vs-dark' : 'light'}
        options={{
          minimap: { enabled: false },
          fontSize: 14,
          lineNumbers: 'on',
          roundedSelection: false,
          scrollBeyondLastLine: false,
          automaticLayout: true,
          tabSize: 2,
          wordWrap: 'on',
        }}
      />
    </div>
  )
}