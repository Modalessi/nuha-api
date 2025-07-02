import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { problemsApi } from '../services/api'
import { ArrowLeft } from 'lucide-react'
import LoadingSpinner from '../components/LoadingSpinner'
import toast from 'react-hot-toast'

export default function EditProblem() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    difficulty: 'EASY',
    tags: [] as string[],
    time_limit: 1,
    memory_limit: 128000,
  })

  useEffect(() => {
    if (id) {
      fetchProblem()
    }
  }, [id])

  const fetchProblem = async () => {
    try {
      setLoading(true)
      const response = await problemsApi.getProblem(id!)
      const problem = response.data
      setFormData({
        title: problem.title,
        description: problem.discription,
        difficulty: problem.difficulty,
        tags: problem.tags || [],
        time_limit: problem.time_limit,
        memory_limit: problem.memory_limit,
      })
    } catch (error) {
      console.error('Error fetching problem:', error)
      toast.error('Problem not found')
      navigate('/admin')
    } finally {
      setLoading(false)
    }
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: name === 'time_limit' || name === 'memory_limit' ? parseFloat(value) : value
    }))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.title.trim() || !formData.description.trim()) {
      toast.error('Title and description are required')
      return
    }

    try {
      setSaving(true)
      await problemsApi.updateProblem(id!, formData)
      toast.success('Problem updated successfully!')
      navigate('/admin')
    } catch (error: any) {
      const message = error.response?.data?.msg || 'Failed to update problem'
      toast.error(message)
    } finally {
      setSaving(false)
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
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <button
          onClick={() => navigate('/admin')}
          className="flex items-center space-x-2 text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="w-4 h-4" />
          <span>Back to Admin</span>
        </button>
        
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Edit Problem</h1>
        <p className="text-gray-600">Update the problem details.</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-8">
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Problem Information</h2>
          
          <div className="space-y-6">
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-2">
                Problem Title *
              </label>
              <input
                type="text"
                id="title"
                name="title"
                value={formData.title}
                onChange={handleInputChange}
                className="input"
                required
              />
            </div>

            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
                Problem Description *
              </label>
              <textarea
                id="description"
                name="description"
                value={formData.description}
                onChange={handleInputChange}
                rows={8}
                className="input"
                required
              />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div>
                <label htmlFor="difficulty" className="block text-sm font-medium text-gray-700 mb-2">
                  Difficulty
                </label>
                <select
                  id="difficulty"
                  name="difficulty"
                  value={formData.difficulty}
                  onChange={handleInputChange}
                  className="input"
                >
                  <option value="EASY">Easy</option>
                  <option value="MEDIUM">Medium</option>
                  <option value="HARD">Hard</option>
                </select>
              </div>

              <div>
                <label htmlFor="time_limit" className="block text-sm font-medium text-gray-700 mb-2">
                  Time Limit (seconds)
                </label>
                <input
                  type="number"
                  id="time_limit"
                  name="time_limit"
                  value={formData.time_limit}
                  onChange={handleInputChange}
                  min="0.1"
                  step="0.1"
                  className="input"
                />
              </div>

              <div>
                <label htmlFor="memory_limit" className="block text-sm font-medium text-gray-700 mb-2">
                  Memory Limit (KB)
                </label>
                <input
                  type="number"
                  id="memory_limit"
                  name="memory_limit"
                  value={formData.memory_limit}
                  onChange={handleInputChange}
                  min="1024"
                  step="1024"
                  className="input"
                />
              </div>
            </div>
          </div>
        </div>

        <div className="flex justify-end space-x-4">
          <button
            type="button"
            onClick={() => navigate('/admin')}
            className="btn-secondary"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={saving}
            className="btn-primary flex items-center space-x-2"
          >
            {saving ? (
              <LoadingSpinner size="sm" />
            ) : (
              <span>Update Problem</span>
            )}
          </button>
        </div>
      </form>
    </div>
  )
}