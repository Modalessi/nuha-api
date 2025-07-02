import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { submissionsApi } from '../services/api'
import { ArrowLeft, Clock, Code, User, CheckCircle, XCircle, AlertCircle } from 'lucide-react'
import LoadingSpinner from '../components/LoadingSpinner'
import { formatDistanceToNow } from 'date-fns'

interface Submission {
  id: string
  problem_id: string
  user_id: string
  status: string
  language: string
  code: string
  created_at: string
}

export default function SubmissionDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [submission, setSubmission] = useState<Submission | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (id) {
      fetchSubmission()
    }
  }, [id])

  const fetchSubmission = async () => {
    try {
      setLoading(true)
      const response = await submissionsApi.getSubmission(id!)
      setSubmission(response.data)
    } catch (error) {
      console.error('Error fetching submission:', error)
      navigate('/submissions')
    } finally {
      setLoading(false)
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'accepted':
        return <CheckCircle className="w-6 h-6 text-success-600" />
      case 'pending':
        return <AlertCircle className="w-6 h-6 text-warning-600" />
      default:
        return <XCircle className="w-6 h-6 text-error-600" />
    }
  }

  const getStatusBadge = (status: string) => {
    switch (status.toLowerCase()) {
      case 'accepted':
        return 'badge-accepted'
      case 'pending':
        return 'badge-pending'
      default:
        return 'badge-wrong'
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  if (!submission) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">Submission not found</h2>
          <p className="text-gray-600 mb-4">The submission you're looking for doesn't exist.</p>
          <button onClick={() => navigate('/submissions')} className="btn-primary">
            Back to Submissions
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
          onClick={() => navigate('/submissions')}
          className="flex items-center space-x-2 text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="w-4 h-4" />
          <span>Back to Submissions</span>
        </button>
        
        <div className="flex items-center space-x-4">
          <h1 className="text-3xl font-bold text-gray-900">Submission Details</h1>
          <div className="flex items-center space-x-2">
            {getStatusIcon(submission.status)}
            <span className={getStatusBadge(submission.status)}>
              {submission.status}
            </span>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Submission Info */}
        <div className="lg:col-span-1 space-y-6">
          <div className="card">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Submission Info</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Submission ID
                </label>
                <p className="text-sm text-gray-900 font-mono">{submission.id}</p>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Problem
                </label>
                <Link
                  to={`/problems/${submission.problem_id}`}
                  className="text-sm text-primary-600 hover:text-primary-900"
                >
                  Problem #{submission.problem_id.slice(0, 8)}
                </Link>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Language
                </label>
                <div className="flex items-center space-x-1">
                  <Code className="w-4 h-4 text-gray-400" />
                  <span className="text-sm text-gray-900">{submission.language}</span>
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Submitted
                </label>
                <div className="flex items-center space-x-1">
                  <Clock className="w-4 h-4 text-gray-400" />
                  <span className="text-sm text-gray-900">
                    {formatDistanceToNow(new Date(submission.created_at), { addSuffix: true })}
                  </span>
                </div>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  User
                </label>
                <div className="flex items-center space-x-1">
                  <User className="w-4 h-4 text-gray-400" />
                  <span className="text-sm text-gray-900">
                    User #{submission.user_id.slice(0, 8)}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Test Results */}
          <div className="card">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Test Results</h2>
            <div className="space-y-3">
              {/* Mock test results - in real app, this would come from your API */}
              {[1, 2, 3, 4, 5].map((testCase) => (
                <div key={testCase} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                  <span className="text-sm font-medium">Test Case {testCase}</span>
                  <div className="flex items-center space-x-2">
                    {submission.status.toLowerCase() === 'accepted' ? (
                      <CheckCircle className="w-4 h-4 text-success-600" />
                    ) : (
                      <XCircle className="w-4 h-4 text-error-600" />
                    )}
                    <span className="text-xs text-gray-500">0.1s</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Code */}
        <div className="lg:col-span-2">
          <div className="card">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Submitted Code</h2>
            <div className="bg-gray-900 rounded-lg overflow-hidden">
              <div className="bg-gray-800 px-4 py-2 border-b border-gray-700">
                <span className="text-sm text-gray-300">{submission.language}</span>
              </div>
              <pre className="p-4 text-sm text-gray-100 overflow-x-auto">
                <code>{submission.code}</code>
              </pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}