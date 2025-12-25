import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { goalsAPI } from '../services/api';
import { ArrowLeft, Plus, Target, Check, Trash2, Loader2, X, TrendingUp, Edit3 } from 'lucide-react';
import './Goals.css';

const Goals = () => {
    const navigate = useNavigate();
    const [goals, setGoals] = useState([]);
    const [stats, setStats] = useState({ total: 0, completed: 0, in_progress: 0 });
    const [loading, setLoading] = useState(true);
    const [showNewGoal, setShowNewGoal] = useState(false);
    const [editingGoalId, setEditingGoalId] = useState(null);
    const [editProgress, setEditProgress] = useState('');
    const [newGoal, setNewGoal] = useState({
        title: '',
        description: '',
        type: 'custom',
        target: '',
        unit: '',
        deadline: ''
    });
    const [submitting, setSubmitting] = useState(false);

    const goalTypes = [
        { id: 'weight', name: 'Berat Badan', icon: '⚖️', unit: 'kg' },
        { id: 'exercise', name: 'Olahraga', icon: '🏃', unit: 'menit/minggu' },
        { id: 'water', name: 'Minum Air', icon: '💧', unit: 'gelas/hari' },
        { id: 'sleep', name: 'Tidur', icon: '😴', unit: 'jam/hari' },
        { id: 'custom', name: 'Lainnya', icon: '🎯', unit: '' }
    ];

    useEffect(() => {
        fetchGoals();
        fetchStats();
    }, []);

    // Lock body scroll when modal is open
    useEffect(() => {
        if (showNewGoal) {
            document.body.style.overflow = 'hidden';
            document.body.classList.add('modal-open');
            document.documentElement.style.overflow = 'hidden';
        } else {
            document.body.style.overflow = '';
            document.body.classList.remove('modal-open');
            document.documentElement.style.overflow = '';
        }
        return () => {
            document.body.style.overflow = '';
            document.body.classList.remove('modal-open');
            document.documentElement.style.overflow = '';
        };
    }, [showNewGoal]);

    const fetchGoals = async () => {
        try {
            const res = await goalsAPI.getAll();
            setGoals(res.data || []);
        } catch (error) {
            console.error('Failed to fetch goals:', error);
        } finally {
            setLoading(false);
        }
    };

    const fetchStats = async () => {
        try {
            const res = await goalsAPI.getStats();
            setStats(res.data);
        } catch (error) {
            console.error('Failed to fetch stats:', error);
        }
    };

    const handleCreateGoal = async (e) => {
        e.preventDefault();
        if (!newGoal.title || !newGoal.target) return;

        setSubmitting(true);
        try {
            // Convert target to number before sending
            const goalData = {
                title: newGoal.title,
                description: newGoal.description || '',
                type: newGoal.type || 'custom',
                target: parseFloat(newGoal.target) || 0,
                unit: newGoal.unit || '',
                deadline: newGoal.deadline || ''
            };
            console.log('Sending goal data:', goalData);
            const response = await goalsAPI.create(goalData);
            console.log('Response:', response);
            setNewGoal({ title: '', description: '', type: 'custom', target: '', unit: '', deadline: '' });
            setShowNewGoal(false);
            fetchGoals();
            fetchStats();
        } catch (error) {
            console.error('Failed to create goal:', error);
            console.error('Error response:', error.response?.data);
            const errorMsg = error.response?.data?.error || 'Gagal menyimpan target. Periksa koneksi backend.';
            alert(errorMsg);
        } finally {
            setSubmitting(false);
        }
    };

    const handleToggleComplete = async (goalId) => {
        try {
            await goalsAPI.toggleComplete(goalId);
            fetchGoals();
            fetchStats();
        } catch (error) {
            console.error('Failed to toggle goal:', error);
        }
    };

    const handleDelete = async (goalId) => {
        if (!confirm('Hapus goal ini?')) return;
        try {
            await goalsAPI.delete(goalId);
            fetchGoals();
            fetchStats();
        } catch (error) {
            console.error('Failed to delete goal:', error);
        }
    };

    const handleStartEditProgress = (goal) => {
        setEditingGoalId(goal.id);
        setEditProgress(goal.current.toString());
    };

    const handleUpdateProgress = async (goalId) => {
        const newValue = parseFloat(editProgress);
        if (isNaN(newValue) || newValue < 0) {
            alert('Masukkan angka yang valid');
            return;
        }
        try {
            await goalsAPI.updateProgress(goalId, newValue);
            setEditingGoalId(null);
            setEditProgress('');
            fetchGoals();
            fetchStats();
        } catch (error) {
            console.error('Failed to update progress:', error);
            alert('Gagal update progress');
        }
    };

    const handleQuickAddProgress = async (goal, increment) => {
        const newValue = (goal.current || 0) + increment;
        try {
            await goalsAPI.updateProgress(goal.id, newValue);
            fetchGoals();
            fetchStats();
        } catch (error) {
            console.error('Failed to add progress:', error);
        }
    };

    const getTypeInfo = (type) => {
        return goalTypes.find(t => t.id === type) || goalTypes[4];
    };

    if (loading) {
        return (
            <div className="goals-loading">
                <Loader2 className="spinner" size={48} />
                <p>Memuat goals...</p>
            </div>
        );
    }

    return (
        <div className="goals-page">
            <header className="goals-header">
                <button className="back-btn" onClick={() => navigate('/dashboard')}>
                    <ArrowLeft size={24} />
                </button>
                <h1>🎯 Target Kesehatan</h1>
                <button className="new-goal-btn" onClick={() => setShowNewGoal(true)}>
                    <Plus size={20} />
                </button>
            </header>

            <div className="goals-container">
                {/* Stats Cards */}
                <div className="goals-stats">
                    <div className="stat-card">
                        <Target size={24} />
                        <div>
                            <span className="stat-value">{stats.total}</span>
                            <span className="stat-label">Total</span>
                        </div>
                    </div>
                    <div className="stat-card completed">
                        <Check size={24} />
                        <div>
                            <span className="stat-value">{stats.completed}</span>
                            <span className="stat-label">Selesai</span>
                        </div>
                    </div>
                    <div className="stat-card progress">
                        <TrendingUp size={24} />
                        <div>
                            <span className="stat-value">{stats.in_progress}</span>
                            <span className="stat-label">Berjalan</span>
                        </div>
                    </div>
                </div>

                {/* Goals List */}
                {goals.length === 0 ? (
                    <div className="empty-state">
                        <Target size={64} />
                        <h3>Belum Ada Target</h3>
                        <p>Buat target kesehatan pertama Anda!</p>
                        <button className="cta-btn" onClick={() => setShowNewGoal(true)}>
                            <Plus size={18} /> Buat Target
                        </button>
                    </div>
                ) : (
                    <div className="goals-list">
                        {goals.map((goal) => {
                            const typeInfo = getTypeInfo(goal.type);
                            return (
                                <div key={goal.id} className={`goal-card ${goal.is_completed ? 'completed' : ''}`}>
                                    <div className="goal-icon">{typeInfo.icon}</div>
                                    <div className="goal-content">
                                        <h3>{goal.title}</h3>
                                        {goal.description && <p>{goal.description}</p>}
                                        <div className="goal-progress">
                                            <div className="progress-bar">
                                                <div
                                                    className="progress-fill"
                                                    style={{ width: `${Math.min(goal.progress || 0, 100)}%` }}
                                                />
                                            </div>
                                            {editingGoalId === goal.id ? (
                                                <div className="progress-edit">
                                                    <input
                                                        type="number"
                                                        value={editProgress}
                                                        onChange={(e) => setEditProgress(e.target.value)}
                                                        autoFocus
                                                        onKeyDown={(e) => {
                                                            if (e.key === 'Enter') handleUpdateProgress(goal.id);
                                                            if (e.key === 'Escape') setEditingGoalId(null);
                                                        }}
                                                    />
                                                    <span>/{goal.target} {goal.unit}</span>
                                                    <button
                                                        className="save-progress-btn"
                                                        onClick={() => handleUpdateProgress(goal.id)}
                                                    >
                                                        <Check size={14} />
                                                    </button>
                                                    <button
                                                        className="cancel-progress-btn"
                                                        onClick={() => setEditingGoalId(null)}
                                                    >
                                                        <X size={14} />
                                                    </button>
                                                </div>
                                            ) : (
                                                <div className="progress-display" onClick={() => handleStartEditProgress(goal)}>
                                                    <span className="progress-text">
                                                        {goal.current || 0}/{goal.target} {goal.unit}
                                                    </span>
                                                    <button
                                                        className="quick-add-btn"
                                                        onClick={(e) => {
                                                            e.stopPropagation();
                                                            handleQuickAddProgress(goal, 1);
                                                        }}
                                                        title="Tambah +1"
                                                    >
                                                        <Plus size={14} />
                                                    </button>
                                                </div>
                                            )}
                                        </div>
                                        {goal.days_left >= 0 && (
                                            <span className="goal-deadline">
                                                {goal.days_left === 0 ? 'Hari ini!' : `${goal.days_left} hari lagi`}
                                            </span>
                                        )}
                                    </div>
                                    <div className="goal-actions">
                                        <button
                                            className={`check-btn ${goal.is_completed ? 'checked' : ''}`}
                                            onClick={() => handleToggleComplete(goal.id)}
                                            title={goal.is_completed ? 'Batalkan selesai' : 'Tandai selesai'}
                                        >
                                            <Check size={18} />
                                        </button>
                                        <button
                                            className="delete-btn"
                                            onClick={() => handleDelete(goal.id)}
                                            title="Hapus target"
                                        >
                                            <Trash2 size={18} />
                                        </button>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>

            {/* New Goal Modal */}
            {showNewGoal && (
                <div className="modal-overlay" onClick={() => setShowNewGoal(false)}>
                    <div className="modal" onClick={(e) => e.stopPropagation()}>
                        <div className="modal-header">
                            <h2>Buat Target Baru</h2>
                            <button className="close-btn" onClick={() => setShowNewGoal(false)}>
                                <X size={24} />
                            </button>
                        </div>
                        <form onSubmit={handleCreateGoal}>
                            <div className="form-group">
                                <label>Jenis Target</label>
                                <div className="type-selector">
                                    {goalTypes.map((type) => (
                                        <button
                                            key={type.id}
                                            type="button"
                                            className={`type-btn ${newGoal.type === type.id ? 'active' : ''}`}
                                            onClick={() => setNewGoal({
                                                ...newGoal,
                                                type: type.id,
                                                unit: type.unit
                                            })}
                                        >
                                            <span>{type.icon}</span>
                                            <span>{type.name}</span>
                                        </button>
                                    ))}
                                </div>
                            </div>
                            <div className="form-group">
                                <label>Judul Target</label>
                                <input
                                    type="text"
                                    placeholder="Contoh: Turunkan berat badan 5kg"
                                    value={newGoal.title}
                                    onChange={(e) => setNewGoal({ ...newGoal, title: e.target.value })}
                                    required
                                />
                            </div>
                            <div className="form-row">
                                <div className="form-group">
                                    <label>Target</label>
                                    <input
                                        type="number"
                                        placeholder="10"
                                        value={newGoal.target}
                                        onChange={(e) => setNewGoal({ ...newGoal, target: e.target.value })}
                                        required
                                    />
                                </div>
                                <div className="form-group">
                                    <label>Satuan</label>
                                    <input
                                        type="text"
                                        placeholder="kg, menit, dll"
                                        value={newGoal.unit}
                                        onChange={(e) => setNewGoal({ ...newGoal, unit: e.target.value })}
                                    />
                                </div>
                            </div>
                            <div className="form-group">
                                <label>Deadline (Opsional)</label>
                                <input
                                    type="date"
                                    value={newGoal.deadline}
                                    onChange={(e) => setNewGoal({ ...newGoal, deadline: e.target.value })}
                                />
                            </div>
                            <button type="submit" className="submit-btn" disabled={submitting}>
                                {submitting ? <Loader2 className="spinner" size={20} /> : <Target size={20} />}
                                {submitting ? 'Menyimpan...' : 'Simpan Target'}
                            </button>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Goals;
