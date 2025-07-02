import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { submissionsApi } from '../services/api'
import { useAuth } from '../contexts/AuthContext'
import { Clock, Code, User, FileText } from 'lucide-react'
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

export default function Submissions() {
  const { user } = useAuth()
  const [submissions, setSubmissions] = useState<Submission[]>([])
  const [loading, setLoading] = useState(true)
  const [currentPage, setCurrentPage] = useState(1)

  useEffect(() => {
    fetchSubmissions()
  }, [currentPage])

  const fetchSubmissions = async () => {
    try {
      setLoading(true)
      const response = await submissionsApi.getSubmissions(currentPage, 20)
      setSubmissions(response.data || [])
    } catch (error) {
      console.error('Error fetching submissions:', error)
    } finally {
      setLoading(false)
    }
  }

  const getStatusBadge = (status: string) => {
    switch (status.toLowerCase()) {
      case 'accepted':
        return 'badge-accepted'
      case 'pending':
        return 'badge-pending'
      case 'wrong answer':
        return 'badge-wrong'
      case 'time limit exceeded':
        return 'badge-wrong'
      case 'compilation error':
        return 'badge-wrong'
      case 'runtime error':
        return 'badge-wrong'
      default:
        return 'badge bg-gray-100 text-gray-800'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'accepted':
        return '✓'
      case 'pending':
        return '⏳'
      default:
        return '✗'
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Submissions</h1>
        <p className="text-gray-600">
          Track your submission history and results.
        </p>
      </div>

      {/* Submissions List */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Problem
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Language
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Submitted
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {submissions.map((submission) => (
                <tr key={submission.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center space-x-2">
                      <span className="text-lg">
                        {getStatusIcon(submission.status)}
                      </span>
                      <span className={getStatusBadge(submission.status)}>
                        {submission.status}
                      </span>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Link
                      to={`/problems/${submission.problem_id}`}
                      className="text-primary-600 hover:text-primary-900 font-medium"
                    >
                      Problem #{submission.problem_id.slice(0, 8)}
                    </Link>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center space-x-1">
                      <Code className="w-4 h-4 text-gray-400" />
                      <span className="text-sm text-gray-900">{submission.language}</span>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center space-x-1">
                      <Clock className="w-4 h-4 text-gray-400" />
                      <span className="text-sm text-gray-500">
                        {formatDistanceToNow(new Date(submission.created_at), { addSuffix: true })}
                      </span>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <Link
                      to={`/submissions/${submission.id}`}
                      className="text-primary-600 hover:text-primary-900 text-sm font-medium"
                    >
                      View Details
                    </Link>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {submissions.length === 0 && !loading && (
        <div className="text-center py-12">
          <div className="text-gray-400 mb-4">
            <FileText className="w-12 h-12 mx-auto" />
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">No submissions yet</h3>
          <p className="text-gray-600 mb-4">Start solving problems to see your submissions here.</p>
          <Link to="/problems" className="btn-primary">
            Browse Problems
          </Link>
        </div>
      )}

      {/* Pagination */}
      {submissions.length > 0 && (
        <div className="flex justify-center mt-8">
          <div className="flex space-x-2">
            <button
              onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
              disabled={currentPage === 1}
              className="btn-secondary disabled:opacity-50"
            >
              Previous
            </button>
            <span className="flex items-center px-4 py-2 text-sm text-gray-700">
              Page {currentPage}
            </span>
            <button
              onClick={() => setCurrentPage(prev => prev + 1)}
              className="btn-secondary"
            >
              Next
            </button>
          </div>
        </div>
      )}
    </div>
  )
}