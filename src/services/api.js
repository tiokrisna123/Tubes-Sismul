import axios from 'axios';

// Menggunakan Link Koyeb Kamu secara langsung
const API_BASE_URL = 'https://bright-swan-1tubes-sismul-cc0e96ef.koyeb.app/api';

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Add token to requests
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Handle auth errors
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

// Auth API
export const authAPI = {
    register: (data) => api.post('/auth/register', data),
    login: (data) => api.post('/auth/login', data),
    getProfile: () => api.get('/auth/me'),
    updateProfile: (data) => api.put('/auth/profile', data),
};

// Health API
export const healthAPI = {
    create: (data) => api.post('/health', data),
    getAll: () => api.get('/health'),
    getLatest: () => api.get('/health/latest'),
    getDashboard: () => api.get('/health/dashboard'),
    getGraph: (period) => api.get(`/health/graph/${period}`),
};

// Symptoms API
export const symptomsAPI = {
    getList: () => api.get('/symptoms/list'),
    log: (data) => api.post('/symptoms', data),
    logBatch: (data) => api.post('/symptoms/batch', data),
    getHistory: () => api.get('/symptoms/history'),
    getStats: () => api.get('/symptoms/stats'),
};

// Family API
export const familyAPI = {
    invite: (data) => api.post('/family/invite', data),
    getMembers: () => api.get('/family/members'),
    getRequests: () => api.get('/family/requests'),
    approve: (id) => api.put(`/family/approve/${id}`),
    reject: (id) => api.put(`/family/reject/${id}`),
    getMemberHealth: (id) => api.get(`/family/${id}/health`),
    remove: (id) => api.delete(`/family/${id}`),
};

// Recommendations API
export const recommendationsAPI = {
    getFood: () => api.get('/recommendations/food'),
    getExercise: () => api.get('/recommendations/exercise'),
    getEmotional: () => api.get('/recommendations/emotional'),
    getDailyMenu: () => api.get('/recommendations/daily-menu'),
};

// Articles API
export const articlesAPI = {
    getAll: (category) => api.get('/articles', { params: { category } }),
    getById: (id) => api.get(`/articles/${id}`),
    getCategories: () => api.get('/articles/categories'),
    search: (q) => api.get('/articles/search', { params: { q } }),
};

// Forum API
export const forumAPI = {
    getPosts: () => api.get('/forum/posts'),
    createPost: (data) => api.post('/forum/posts', data),
    getPost: (id) => api.get(`/forum/posts/${id}`),
    deletePost: (id) => api.delete(`/forum/posts/${id}`),
    addComment: (id, data) => api.post(`/forum/posts/${id}/comments`, data),
    toggleLike: (id) => api.post(`/forum/posts/${id}/like`),
};

// Water Tracker API
export const waterAPI = {
    get: () => api.get('/water'),
    addGlass: () => api.post('/water/add'),
    removeGlass: () => api.post('/water/remove'),
    updateGoal: (goal) => api.put('/water/goal', { goal }),
    getHistory: () => api.get('/water/history'),
};

// Reminders API
export const remindersAPI = {
    getAll: () => api.get('/reminders'),
    create: (data) => api.post('/reminders', data),
    update: (id, data) => api.put(`/reminders/${id}`, data),
    delete: (id) => api.delete(`/reminders/${id}`),
    toggle: (id) => api.put(`/reminders/${id}/toggle`),
};

// Goals API (Saya lihat ada file Goals.jsx yang pakai goalsAPI tapi belum ada di api.js kamu sebelumnya, saya tambahkan biar aman)
export const goalsAPI = {
    getAll: () => api.get('/goals'),
    create: (data) => api.post('/goals', data),
    updateProgress: (id, current) => api.put(`/goals/${id}/progress`, { current }),
    toggleComplete: (id) => api.put(`/goals/${id}/toggle`),
    delete: (id) => api.delete(`/goals/${id}`),
    getStats: () => api.get('/goals/stats'),
};

export default api;