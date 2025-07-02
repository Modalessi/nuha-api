import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { problemsApi } from '../services/api'
import { ArrowLeft, Plus, Trash2 } from 'lucide-react'
import LoadingSpinner from '../components/LoadingSpinner'
import toast from 'react-hot-toast'

interface TestCase {
  stdin: string
  expected_output: string
}

export default function CreateProblem() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    difficulty: 'EASY',
    tags: [] as string[],
    timelimit: 1,
    memorylimit: 128000,
  })
  const [testCases, setTestCases] = useState<TestCase[]>([
    { stdin: '', expected_output: '' }
  ])
  const [tagInput, setTagInput] = useState('')

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: name === 'timelimit' || name === 'memorylimit' ? parseFloat(value) : value
    }))
  }

  const addTag = () => {
    if (tagInput.trim() && !formData.tags.includes(tagInput.trim())) {
      setFormData(prev => ({
        ...prev,
        tags: [...prev.tags, tagInput.trim()]
      }))
      setTagInput('')
    }
  }

  const removeTag = (tagToRemove: string) => {
    setFormData(prev => ({
      ...prev,
      tags: prev.tags.filter(tag => tag !== tagToRemove)
    }))
  }

  const addTestCase = () => {
    setTestCases(prev => [...prev, { stdin: '', expected_output: '' }])
  }

  const removeTestCase = (index: number) => {
    if (testCases.length > 1) {
      setTestCases(prev => prev.filter((_, i) => i !== index))
    }
  }

  const updateTestCase = (index: number, field: keyof TestCase, value: string) => {
    setTestCases(prev => prev.map((tc, i) => 
      i === index ? { ...tc, [field]: value } : tc
    ))
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!formData.title.trim() || !formData.description.trim()) {
      toast.error('Title and description are required')
      return
    }

    if (testCases.some(tc => !tc.stdin.trim() || !tc.expected_output.trim())) {
      toast.error('All test cases must have input and output')
      return
    }

    try {
      setLoading(true)
      
      // Create the problem
      const problemResponse = await problemsApi.createProblem(formData)
      const problemId = problemResponse.data.ID
      
      // Add test cases
      await problemsApi.addTestCases(problemId, testCases)
      
      toast.success('Problem created successfully!')
      navigate('/admin')
    } catch (error: any) {
      const message = error.response?.data?.msg || 'Failed to create problem'
      toast.error(message)
    } finally {
      setLoading(false)
    }
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
        
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Create New Problem</h1>
        <p className="text-gray-600">Add a new competitive programming problem to your platform.</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-8">
        {/* Basic Information */}
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-6">Basic Information</h2>
          
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
                placeholder="Enter problem title"
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
                placeholder="Describe the problem, input/output format, constraints, etc."
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
                <label htmlFor="timelimit" className="block text-sm font-medium text-gray-700 mb-2">
                  Time Limit (seconds)
                </label>
                <input
                  type="number"
                  id="timelimit"
                  name="timelimit"
                  value={formData.timelimit}
                  onChange={handleInputChange}
                  min="0.1"
                  step="0.1"
                  className="input"
                />
              </div>

              <div>
                <label htmlFor="memorylimit" className="block text-sm font-medium text-gray-700 mb-2">
                  Memory Limit (KB)
                </label>
                <input
                  type="number"
                  id="memorylimit"
                  name="memorylimit"
                  value={formData.memorylimit}
                  onChange={handleInputChange}
                  min="1024"
                  step="1024"
                  className="input"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Tags
              </label>
              <div className="flex space-x-2 mb-2">
                <input
                  type="text"
                  value={tagInput}
                  onChange={(e) => setTagInput(e.target.value)}
                  onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
                  className="input flex-1"
                  placeholder="Add a tag"
                />
                <button
                  type="button"
                  onClick={addTag}
                  className="btn-secondary"
                >
                  Add
                </button>
              </div>
              <div className="flex flex-wrap gap-2">
                {formData.tags.map((tag, index) => (
                  <span
                    key={index}
                    className="inline-flex items-center px-3 py-1 bg-primary-100 text-primary-800 text-sm rounded-full"
                  >
                    {tag}
                    <button
                      type="button"
                      onClick={() => removeTag(tag)}
                      className="ml-2 text-primary-600 hover:text-primary-800"
                    >
                      ×
                    </button>
                  </span>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Test Cases */}
        <div className="card">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-xl font-semibold text-gray-900">Test Cases</h2>
            <button
              type="button"
              onClick={addTestCase}
              className="btn-secondary flex items-center space-x-2"
            >
              <Plus className="w-4 h-4" />
              <span>Add Test Case</span>
            </button>
          </div>

          <div className="space-y-6">
            {testCases.map((testCase, index) => (
              <div key={index} className="border border-gray-200 rounded-lg p-4">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-medium text-gray-900">Test Case {index + 1}</h3>
                  {testCases.length > 1 && (
                    <button
                      type="button"
                      onClick={() => removeTestCase(index)}
                      className="text-red-600 hover:text-red-800"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  )}
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Input
                    </label>
                    <textarea
                      value={testCase.stdin}
                      onChange={(e) => updateTestCase(index, 'stdin', e.target.value)}
                      rows={4}
                      className="input"
                      placeholder="Enter test input"
                    />
                  </div>
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      Expected Output
                    </label>
                    <textarea
                      value={testCase.expected_output}
                      onChange={(e) => updateTestCase(index, 'expected_output', e.target.value)}
                      rows={4}
                      className="input"
                      placeholder="Enter expected output"
                    />
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Submit Button */}
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
            disabled={loading}
            className="btn-primary flex items-center space-x-2"
          >
            {loading ? (
              <LoadingSpinner size="sm" />
            ) : (
              <span>Create Problem</span>
            )}
          </button>
        </div>
      </form>
    </div>
  )
}