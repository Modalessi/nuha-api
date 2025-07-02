import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Edit, Trash2, Users, FileText, Trophy, Settings } from 'lucide-react'

export default function AdminDashboard() {
  const [activeTab, setActiveTab] = useState('overview')

  const stats = [
    { label: 'Total Problems', value: '156', icon: FileText, color: 'text-blue-600' },
    { label: 'Total Users', value: '2,341', icon: Users, color: 'text-green-600' },
    { label: 'Total Submissions', value: '12,456', icon: Trophy, color: 'text-purple-600' },
    { label: 'Active Contests', value: '3', icon: Settings, color: 'text-orange-600' },
  ]

  const recentProblems = [
    { id: '1', title: 'Two Sum', difficulty: 'Easy', submissions: 234 },
    { id: '2', title: 'Binary Search Tree', difficulty: 'Medium', submissions: 156 },
    { id: '3', title: 'Graph Traversal', difficulty: 'Hard', submissions: 89 },
  ]

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">Admin Dashboard</h1>
        <p className="text-gray-600">Manage your competitive programming platform.</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
        {/* Sidebar */}
        <div className="lg:col-span-1">
          <div className="card">
            <nav className="space-y-2">
              <button
                onClick={() => setActiveTab('overview')}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === 'overview'
                    ? 'bg-primary-50 text-primary-600'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                }`}
              >
                Overview
              </button>
              <button
                onClick={() => setActiveTab('problems')}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === 'problems'
                    ? 'bg-primary-50 text-primary-600'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                }`}
              >
                Problems
              </button>
              <button
                onClick={() => setActiveTab('users')}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === 'users'
                    ? 'bg-primary-50 text-primary-600'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                }`}
              >
                Users
              </button>
              <button
                onClick={() => setActiveTab('submissions')}
                className={`w-full text-left px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  activeTab === 'submissions'
                    ? 'bg-primary-50 text-primary-600'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                }`}
              >
                Submissions
              </button>
            </nav>
          </div>
        </div>

        {/* Main Content */}
        <div className="lg:col-span-3">
          {activeTab === 'overview' && (
            <div className="space-y-6">
              {/* Stats */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {stats.map((stat, index) => (
                  <div key={index} className="card text-center">
                    <div className={`w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center mx-auto mb-3`}>
                      <stat.icon className={`w-6 h-6 ${stat.color}`} />
                    </div>
                    <div className="text-2xl font-bold text-gray-900 mb-1">{stat.value}</div>
                    <div className="text-sm text-gray-600">{stat.label}</div>
                  </div>
                ))}
              </div>

              {/* Quick Actions */}
              <div className="card">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <Link
                    to="/admin/problems/create"
                    className="flex items-center space-x-3 p-4 bg-primary-50 rounded-lg hover:bg-primary-100 transition-colors"
                  >
                    <Plus className="w-6 h-6 text-primary-600" />
                    <span className="font-medium text-primary-900">Create Problem</span>
                  </Link>
                  <button className="flex items-center space-x-3 p-4 bg-green-50 rounded-lg hover:bg-green-100 transition-colors">
                    <Users className="w-6 h-6 text-green-600" />
                    <span className="font-medium text-green-900">Manage Users</span>
                  </button>
                  <button className="flex items-center space-x-3 p-4 bg-purple-50 rounded-lg hover:bg-purple-100 transition-colors">
                    <Trophy className="w-6 h-6 text-purple-600" />
                    <span className="font-medium text-purple-900">Create Contest</span>
                  </button>
                </div>
              </div>

              {/* Recent Problems */}
              <div className="card">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-semibold text-gray-900">Recent Problems</h3>
                  <Link to="/admin/problems/create" className="btn-primary">
                    <Plus className="w-4 h-4 mr-2" />
                    Add Problem
                  </Link>
                </div>
                <div className="space-y-3">
                  {recentProblems.map((problem) => (
                    <div key={problem.id} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                      <div>
                        <div className="font-medium text-gray-900">{problem.title}</div>
                        <div className="text-sm text-gray-500">
                          {problem.difficulty} • {problem.submissions} submissions
                        </div>
                      </div>
                      <div className="flex items-center space-x-2">
                        <button className="p-2 text-gray-400 hover:text-gray-600">
                          <Edit className="w-4 h-4" />
                        </button>
                        <button className="p-2 text-gray-400 hover:text-red-600">
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          {activeTab === 'problems' && (
            <div className="card">
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-semibold text-gray-900">Problem Management</h3>
                <Link to="/admin/problems/create" className="btn-primary">
                  <Plus className="w-4 h-4 mr-2" />
                  Create Problem
                </Link>
              </div>
              <p className="text-gray-600">Manage all problems in your platform.</p>
            </div>
          )}

          {activeTab === 'users' && (
            <div className="card">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">User Management</h3>
              <p className="text-gray-600">View and manage platform users.</p>
            </div>
          )}

          {activeTab === 'submissions' && (
            <div className="card">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Submission Management</h3>
              <p className="text-gray-600">Monitor and manage user submissions.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}