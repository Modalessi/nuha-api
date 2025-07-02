import axios from 'axios'

// Configure the base URL for your Go backend
export const API_BASE_URL = 'http://localhost:8080' // Update this to your backend URL

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// API service functions
export const problemsApi = {
  getProblems: (page = 1, perPage = 20) => 
    api.get(`/problem?page=${page}&per_page=${perPage}`),
  
  getProblem: (id: string) => 
    api.get(`/problem?problem_id=${id}`),
  
  createProblem: (data: any) => 
    api.post('/problem', data),
  
  updateProblem: (id: string, data: any) => 
    api.put(`/problem?problem_id=${id}`, data),
  
  deleteProblem: (id: string) => 
    api.delete(`/problem?problem_id=${id}`),
  
  addTestCases: (problemId: string, testcases: any) => 
    api.post(`/testcase?problem_id=${problemId}`, testcases),
}

export const submissionsApi = {
  submitSolution: (problemId: string, data: { language: number; code: string }) =>
    api.post(`/submit?problem_id=${problemId}`, data),
  
  getSubmissions: (page = 1, perPage = 20) =>
    api.get(`/submit?page=${page}&per_page=${perPage}`),
  
  getSubmission: (id: string) =>
    api.get(`/submit?submission_id=${id}`),
  
  getUserSubmissions: (userId: string, page = 1, perPage = 20) =>
    api.get(`/submit?user_id=${userId}&page=${page}&per_page=${perPage}`),
  
  getProblemSubmissions: (problemId: string, page = 1, perPage = 20) =>
    api.get(`/submit?problem_id=${problemId}&page=${page}&per_page=${perPage}`),
}

export const authApi = {
  login: (email: string, password: string) =>
    api.post('/login', { email, password }),
  
  register: (data: any) =>
    api.post('/register', data),
  
  logout: () =>
    api.post('/logout'),
  
  verifyEmail: (token: string) =>
    api.get(`/verify/${token}`),
}