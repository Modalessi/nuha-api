import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { problemsApi, submissionsApi } from '../services/api'
import { useAuth } from '../contexts/AuthContext'
import { Clock, MemoryStick, Tag, Play, ArrowLeft } from 'lucide-react'
import CodeEditor from '../components/CodeEditor'
import LoadingSpinner from '../components/LoadingSpinner'
import toast from 'react-hot-toast'

interface Problem {
  id: string
  title: string
  difficulty: string
  discription: string
  tags: string[]
  time_limit: number
  memory_limit: number
}

const LANGUAGE_OPTIONS = [
  { value: '71', label: 'Python 3.8' },
  { value: '92', label: 'Python 3.11' },
  { value: '100', label: 'Python 3.12' },
  { value: '62', label: 'Java' },
  { value: '54', label: 'C++ (GCC 9.2.0)' },
  { value: '53', label: 'C++ (GCC 8.3.0)' },
  { value: '50', label: 'C (GCC 9.2.0)' },
  { value: '63', label: 'JavaScript' },
  { value: '74', label: 'TypeScript' },
]

export default function ProblemDetail() {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuth()
  const navigate = useNavigate()
  
  const [problem, setProblem] = useState<Problem | null>(null)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [code, setCode] = useState('')
  const [selectedLanguage, setSelectedLanguage] = useState('71')

  useEffect(() => {
    if (id) {
      fetchProblem()
    }
  }, [id])

  const fetchProblem = async () => {
    try {
      setLoading(true)
      const response = await problemsApi.getProblem(id!)
      setProblem(response.data)
    } catch (error) {
      console.error('Error fetching problem:', error)
      toast.error('Problem not found')
      navigate('/problems')
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async () => {
    if (!user) {
      toast.error('Please login to submit solutions')
      navigate('/login')
      return
    }

    if (!code.trim()) {
      toast.error('Please write some code before submitting')
      return
    }

    try {
      setSubmitting(true)
      const response = await submissionsApi.submitSolution(id!, {
        language: parseInt(selectedLanguage),
        code: code
      })
      
      toast.success('Solution submitted successfully!')
      navigate(`/submissions/${response.data.submission_id}`)
    } catch (error: any) {
      const message = error.response?.data?.msg || 'Failed to submit solution'
      toast.error(message)
    } finally {
      setSubmitting(false)
    }
  }

  const getDifficultyBadge = (difficulty: string) => {
    switch (difficulty.toLowerCase()) {
      case 'easy':
        return 'badge-easy'
      case 'medium':
        return 'badge-medium'
      case 'hard':
        return 'badge-hard'
      default:
        return 'badge bg-gray-100 text-gray-800'
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  if (!problem) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Problem not found</h2>
          <p className="text-gray-600 mb-4">The problem you're looking for doesn't exist.</p>
          <button onClick={() => navigate('/problems')} className="btn-primary">
            Back to Problems
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <button
          onClick={() => navigate('/problems')}
          className="flex items-center space-x-2 text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="w-4 h-4" />
          <span>Back to Problems</span>
        </button>
        
        <div className="flex items-start justify-between">
          <div>
            <div className="flex items-center space-x-3 mb-2">
              <h1 className="text-3xl font-bold text-gray-900">{problem.title}</h1>
              <span className={getDifficultyBadge(problem.difficulty)}>
                {problem.difficulty}
              </span>
            </div>
            
            <div className="flex items-center space-x-6 text-sm text-gray-500 mb-4">
              <div className="flex items-center space-x-1">
                <Clock className="w-4 h-4" />
                <span>Time Limit: {problem.time_limit}s</span>
              </div>
              <div className="flex items-center space-x-1">
                <MemoryStick className="w-4 h-4" />
                <span>Memory Limit: {Math.round(problem.memory_limit / 1024)}MB</span>
              </div>
            </div>
            
            <div className="flex items-center space-x-2">
              <Tag className="w-4 h-4 text-gray-400" />
              <div className="flex flex-wrap gap-2">
                {problem.tags.map((tag, index) => (
                  <span
                    key={index}
                    className="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-md"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Problem Description */}
        <div className="space-y-6">
          <div className="card">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Problem Description</h2>
            <div className="prose prose-sm max-w-none">
              <div className="whitespace-pre-wrap text-gray-700">
                {problem.discription || 'No description available.'}
              </div>
            </div>
          </div>

          {/* Sample Input/Output would go here if available */}
          <div className="card">
            <h3 className="text-lg font-semibold text-gray-900 mb-4">Sample Input/Output</h3>
            <div className="space-y-4">
              <div>
                <h4 className="font-medium text-gray-900 mb-2">Input:</h4>
                <pre className="bg-gray-50 p-3 rounded-lg text-sm">
                  {/* Sample input would come from your API */}
                  5 3
                </pre>
              </div>
              <div>
                <h4 className="font-medium text-gray-900 mb-2">Output:</h4>
                <pre className="bg-gray-50 p-3 rounded-lg text-sm">
                  {/* Sample output would come from your API */}
                  8
                </pre>
              </div>
            </div>
          </div>
        </div>

        {/* Code Editor */}
        <div className="space-y-6">
          <div className="card">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-semibold text-gray-900">Solution</h2>
              <select
                value={selectedLanguage}
                onChange={(e) => setSelectedLanguage(e.target.value)}
                className="input w-auto"
              >
                {LANGUAGE_OPTIONS.map((lang) => (
                  <option key={lang.value} value={lang.value}>
                    {lang.label}
                  </option>
                ))}
              </select>
            </div>
            
            <CodeEditor
              value={code}
              onChange={setCode}
              language={selectedLanguage}
              height="500px"
            />
            
            <div className="mt-4 flex justify-end">
              <button
                onClick={handleSubmit}
                disabled={submitting || !user}
                className="btn-primary flex items-center space-x-2"
              >
                {submitting ? (
                  <LoadingSpinner size="sm" />
                ) : (
                  <Play className="w-4 h-4" />
                )}
                <span>{submitting ? 'Submitting...' : 'Submit Solution'}</span>
              </button>
            </div>
            
            {!user && (
              <p className="text-sm text-gray-500 mt-2 text-center">
                Please <button onClick={() => navigate('/login')} className="text-primary-600 hover:underline">login</button> to submit solutions
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}