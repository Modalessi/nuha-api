import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from './contexts/AuthContext'
import Layout from './components/Layout'
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import Problems from './pages/Problems'
import ProblemDetail from './pages/ProblemDetail'
import Submissions from './pages/Submissions'
import SubmissionDetail from './pages/SubmissionDetail'
import Profile from './pages/Profile'
import AdminDashboard from './pages/AdminDashboard'
import CreateProblem from './pages/CreateProblem'
import EditProblem from './pages/EditProblem'
import LoadingSpinner from './components/LoadingSpinner'

function App() {
  const { user, loading } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="lg" />
      </div>
    )
  }

  return (
    <Routes>
      <Route path="/login" element={!user ? <Login /> : <Navigate to="/" />} />
      <Route path="/register" element={!user ? <Register /> : <Navigate to="/" />} />
      
      <Route path="/" element={<Layout />}>
        <Route index element={<Home />} />
        <Route path="problems" element={<Problems />} />
        <Route path="problems/:id" element={<ProblemDetail />} />
        <Route path="submissions" element={user ? <Submissions /> : <Navigate to="/login" />} />
        <Route path="submissions/:id" element={user ? <SubmissionDetail /> : <Navigate to="/login" />} />
        <Route path="profile" element={user ? <Profile /> : <Navigate to="/login" />} />
        
        {/* Admin Routes */}
        <Route path="admin" element={user ? <AdminDashboard /> : <Navigate to="/login" />} />
        <Route path="admin/problems/create" element={user ? <CreateProblem /> : <Navigate to="/login" />} />
        <Route path="admin/problems/:id/edit" element={user ? <EditProblem /> : <Navigate to="/login" />} />
      </Route>
    </Routes>
  )
}

export default App