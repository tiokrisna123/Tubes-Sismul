import { useState, useEffect, useCallback } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import { healthAPI, symptomsAPI, remindersAPI } from '../services/api';
import {
    Heart, Activity, Utensils, Dumbbell, Users, Brain,
    TrendingUp, AlertCircle, Calendar, LogOut, Plus, ChevronRight, X,
    Droplets, Coffee, Moon, Sun, Wind, Bell, Clock, Smile, Frown, Meh,
    AlertTriangle, CheckCircle, History, BarChart3, Edit3, Trash2
} from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, AreaChart, Area } from 'recharts';
import './Dashboard.css';

const Dashboard = () => {
    const { user, logout } = useAuth();
    const { isDark, toggleTheme } = useTheme();
    const navigate = useNavigate();
    const [dashboard, setDashboard] = useState(null);
    const [graphData, setGraphData] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showSymptomsModal, setShowSymptomsModal] = useState(false);
    const [showHistoryModal, setShowHistoryModal] = useState(false);
    const [showRemindersPanel, setShowRemindersPanel] = useState(false);
    const [graphPeriod, setGraphPeriod] = useState('week');
    const [symptomHistory, setSymptomHistory] = useState([]);
    const [healthAlerts, setHealthAlerts] = useState([]);

    // Default reminder templates
    const defaultReminders = [
        { id: 'default-1', type: 'water', label: 'Minum Air Pagi', time: '07:00', is_active: true },
        { id: 'default-2', type: 'meal', label: 'Sarapan Sehat', time: '07:30', is_active: true },
        { id: 'default-3', type: 'exercise', label: 'Olahraga Pagi', time: '06:00', is_active: true },
        { id: 'default-4', type: 'meditation', label: 'Meditasi', time: '06:30', is_active: true },
        { id: 'default-5', type: 'water', label: 'Minum Air Siang', time: '12:00', is_active: true },
        { id: 'default-6', type: 'meal', label: 'Makan Siang', time: '12:30', is_active: true },
        { id: 'default-7', type: 'rest', label: 'Istirahat Siang', time: '13:00', is_active: true },
        { id: 'default-8', type: 'water', label: 'Minum Air Sore', time: '16:00', is_active: true },
        { id: 'default-9', type: 'meal', label: 'Makan Malam', time: '19:00', is_active: true },
        { id: 'default-10', type: 'rest', label: 'Persiapan Tidur', time: '21:00', is_active: true },
    ];

    // Health Reminders State - start with default templates
    const [reminders, setReminders] = useState(defaultReminders);
    const [showReminderModal, setShowReminderModal] = useState(false);
    const [editingReminder, setEditingReminder] = useState(null);
    const [reminderForm, setReminderForm] = useState({
        type: 'water',
        label: '',
        time: '08:00'
    });
    const [reminderLoading, setReminderLoading] = useState(false);

    useEffect(() => {
        loadDashboard();
        loadSymptomHistory();
        loadReminders();
    }, []);

    useEffect(() => {
        loadGraphData(graphPeriod);
    }, [graphPeriod]);

    const loadDashboard = async () => {
        try {
            const dashRes = await healthAPI.getDashboard();
            setDashboard(dashRes.data.data);

            // Analyze patterns for alerts
            analyzeHealthPatterns(dashRes.data.data);
        } catch (err) {
            console.error('Failed to load dashboard:', err);
        } finally {
            setLoading(false);
        }
    };

    const loadGraphData = async (period) => {
        try {
            const graphRes = await healthAPI.getGraph(period);
            setGraphData(graphRes.data.data || []);
        } catch (err) {
            console.error('Failed to load graph data:', err);
        }
    };

    const loadSymptomHistory = async () => {
        try {
            const res = await symptomsAPI.getHistory();
            setSymptomHistory(res.data.data || []);
        } catch (err) {
            console.error('Failed to load symptom history:', err);
        }
    };

    // Analyze health patterns for priority alerts
    const analyzeHealthPatterns = (data) => {
        const alerts = [];
        const symptoms = data?.recent_symptoms || [];

        // Check for recurring stress
        const stressSymptoms = symptoms.filter(s =>
            s.symptom_name?.toLowerCase().includes('stres') ||
            s.symptom_name?.toLowerCase().includes('cemas')
        );
        if (stressSymptoms.length >= 3) {
            alerts.push({
                type: 'warning',
                title: 'Pola Stres Terdeteksi',
                message: 'Anda mengalami stres berulang. Pertimbangkan konsultasi dengan psikolog.',
                priority: 'high'
            });
        }

        // Check for sleep issues
        const sleepSymptoms = symptoms.filter(s =>
            s.symptom_name?.toLowerCase().includes('tidur') ||
            s.symptom_name?.toLowerCase().includes('insomnia')
        );
        if (sleepSymptoms.length >= 2) {
            alerts.push({
                type: 'warning',
                title: 'Gangguan Tidur Berkepanjangan',
                message: 'Pola tidur Anda terganggu. Hindari kafein malam hari dan coba teknik relaksasi.',
                priority: 'high'
            });
        }

        // Check for physical symptoms
        const physicalSymptoms = symptoms.filter(s => s.symptom_type === 'physical');
        if (physicalSymptoms.length >= 4) {
            alerts.push({
                type: 'danger',
                title: 'Banyak Gejala Fisik',
                message: 'Anda memiliki beberapa gejala fisik. Sangat disarankan untuk konsultasi ke dokter.',
                priority: 'critical'
            });
        }

        // Health score alert
        const healthScore = data?.health_score || 0;
        if (healthScore < 60) {
            alerts.push({
                type: 'info',
                title: 'Skor Kesehatan Perlu Perhatian',
                message: 'Skor kesehatan Anda di bawah optimal. Fokus pada pola makan, olahraga, dan istirahat.',
                priority: 'medium'
            });
        }

        setHealthAlerts(alerts);
    };

    // Reminders CRUD functions
    const loadReminders = async () => {
        try {
            const res = await remindersAPI.getAll();
            // Only update if API returns data, otherwise keep default templates
            if (res.data.data && res.data.data.length > 0) {
                setReminders(res.data.data);
            }
        } catch (err) {
            console.error('Failed to load reminders, using default templates:', err);
            // Keep using default templates if API fails
        }
    };

    const toggleReminder = async (id) => {
        // Check if this is a default template (not from API)
        const isDefaultTemplate = String(id).startsWith('default-');

        if (isDefaultTemplate) {
            // Toggle locally for default templates
            setReminders(reminders.map(r =>
                r.id === id ? { ...r, is_active: !r.is_active } : r
            ));
        } else {
            try {
                await remindersAPI.toggle(id);
                setReminders(reminders.map(r =>
                    r.id === id ? { ...r, is_active: !r.is_active } : r
                ));
            } catch (err) {
                console.error('Failed to toggle reminder:', err);
            }
        }
    };

    const getReminderIcon = (type) => {
        const icons = {
            water: Droplets,
            meal: Coffee,
            exercise: Dumbbell,
            meditation: Wind,
            rest: Moon,
            custom: Bell
        };
        return icons[type] || Bell;
    };

    const openAddReminderModal = () => {
        setEditingReminder(null);
        setReminderForm({ type: 'water', label: '', time: '08:00' });
        setShowReminderModal(true);
    };

    const openEditReminderModal = (reminder) => {
        setEditingReminder(reminder);
        setReminderForm({
            type: reminder.type,
            label: reminder.label,
            time: reminder.time
        });
        setShowReminderModal(true);
    };

    const closeReminderModal = () => {
        setShowReminderModal(false);
        setEditingReminder(null);
        setReminderForm({ type: 'water', label: '', time: '08:00' });
    };

    const handleSaveReminder = async (e) => {
        e.preventDefault();
        if (!reminderForm.label.trim()) return;

        setReminderLoading(true);
        try {
            if (editingReminder) {
                await remindersAPI.update(editingReminder.id, reminderForm);
            } else {
                await remindersAPI.create(reminderForm);
            }
            await loadReminders();
            closeReminderModal();
        } catch (err) {
            console.error('Failed to save reminder:', err);
        } finally {
            setReminderLoading(false);
        }
    };

    const handleDeleteReminder = async (id) => {
        if (!window.confirm('Hapus pengingat ini?')) return;

        const isDefaultTemplate = String(id).startsWith('default-');

        if (isDefaultTemplate) {
            // Delete locally for default templates
            setReminders(reminders.filter(r => r.id !== id));
        } else {
            try {
                await remindersAPI.delete(id);
                setReminders(reminders.filter(r => r.id !== id));
            } catch (err) {
                console.error('Failed to delete reminder:', err);
            }
        }
    };

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    const getBMIColor = (category) => {
        switch (category) {
            case 'Underweight': return '#f59e0b';
            case 'Normal': return '#10b981';
            case 'Overweight': return '#f59e0b';
            case 'Obese': return '#ef4444';
            default: return '#64748b';
        }
    };

    const getScoreColor = (score) => {
        if (score >= 80) return '#10b981';
        if (score >= 60) return '#f59e0b';
        return '#ef4444';
    };

    const getMoodIcon = (mood) => {
        switch (mood) {
            case 'happy': return <Smile className="mood-icon happy" />;
            case 'sad': return <Frown className="mood-icon sad" />;
            case 'stressed':
            case 'anxious': return <Heart className="mood-icon stressed" />;
            default: return <Meh className="mood-icon neutral" />;
        }
    };

    const getMoodLabel = (mood) => {
        switch (mood) {
            case 'happy': return 'Bahagia';
            case 'sad': return 'Sedih';
            case 'stressed': return 'Stres';
            case 'anxious': return 'Cemas';
            case 'neutral': return 'Netral';
            default: return 'Tidak diketahui';
        }
    };

    const getCompletedReminders = () => reminders.filter(r => !r.is_active).length;

    if (loading) {
        return (
            <div className="dashboard-loading">
                <Heart className="loading-icon" />
                <p>Memuat dashboard...</p>
            </div>
        );
    }

    const bmi = dashboard?.latest_health?.bmi || 0;
    const bmiCategory = dashboard?.bmi_category || 'Unknown';
    const healthScore = dashboard?.health_score || 0;
    const emotionalState = dashboard?.latest_health?.emotional_state || 'neutral';

    return (
        <div className="dashboard">
            {/* Animated Background Elements */}
            <div className="dashboard-bg-elements">
                {/* Floating Particles */}
                <div className="dash-particle"></div>
                <div className="dash-particle"></div>
                <div className="dash-particle"></div>
                <div className="dash-particle"></div>
                <div className="dash-particle"></div>
                <div className="dash-particle"></div>
                <div className="dash-particle"></div>
                <div className="dash-particle"></div>
                <div className="dash-particle"></div>
                <div className="dash-particle"></div>
                <div className="dash-particle"></div>
                <div className="dash-particle"></div>

                {/* Glowing Orbs */}
                <div className="dash-glow-orb"></div>
                <div className="dash-glow-orb"></div>
                <div className="dash-glow-orb"></div>
                <div className="dash-glow-orb"></div>
                <div className="dash-glow-orb"></div>
                <div className="dash-glow-orb"></div>

                {/* Twinkling Stars */}
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>
                <div className="dash-star"></div>

                {/* Floating Hexagons */}
                <div className="dash-hexagon"></div>
                <div className="dash-hexagon"></div>
                <div className="dash-hexagon"></div>
                <div className="dash-hexagon"></div>

                {/* Connecting Lines */}
                <div className="dash-line"></div>
                <div className="dash-line"></div>
                <div className="dash-line"></div>

                {/* Floating Health Bubbles */}
                <div className="dash-bubble"></div>
                <div className="dash-bubble"></div>
                <div className="dash-bubble"></div>
                <div className="dash-bubble"></div>
                <div className="dash-bubble"></div>
                <div className="dash-bubble"></div>
                <div className="dash-bubble"></div>
                <div className="dash-bubble"></div>

                {/* Pulse Rings */}
                <div className="dash-pulse-ring"></div>
                <div className="dash-pulse-ring"></div>
                <div className="dash-pulse-ring"></div>
            </div>

            <nav className="dashboard-nav" role="navigation" aria-label="Main navigation">
                <div className="nav-brand">
                    <Heart className="nav-logo" aria-hidden="true" />
                    <span>Live for Health</span>
                </div>
                <div className="nav-user">
                    <button
                        className="theme-toggle-btn"
                        onClick={toggleTheme}
                        aria-label={isDark ? "Switch to light mode" : "Switch to dark mode"}
                    >
                        {isDark ? <Sun size={18} /> : <Moon size={18} />}
                    </button>
                    <button
                        className="reminder-btn"
                        onClick={() => setShowRemindersPanel(!showRemindersPanel)}
                        aria-label="Toggle reminders panel"
                    >
                        <Bell size={18} aria-hidden="true" />
                        <span className="reminder-badge" aria-label={`${reminders.length - getCompletedReminders()} active reminders`}>
                            {reminders.length - getCompletedReminders()}
                        </span>
                    </button>
                    <Link to="/profile" className="profile-link" aria-label="Open profile settings">
                        <span>Halo, {user?.name || 'User'}</span>
                    </Link>
                    <button className="logout-btn" onClick={handleLogout} aria-label="Logout">
                        <LogOut size={18} aria-hidden="true" />
                    </button>
                </div>
            </nav>

            {/* Health Alerts Banner */}
            {healthAlerts.length > 0 && (
                <div className="health-alerts-banner">
                    {healthAlerts.map((alert, idx) => (
                        <div key={idx} className={`alert-item alert-${alert.type}`}>
                            <AlertTriangle size={18} />
                            <div className="alert-content">
                                <strong>{alert.title}</strong>
                                <p>{alert.message}</p>
                            </div>
                            <button onClick={() => setHealthAlerts(healthAlerts.filter((_, i) => i !== idx))}>
                                <X size={16} />
                            </button>
                        </div>
                    ))}
                </div>
            )}

            <div className="dashboard-content">
                <header className="dashboard-header">
                    <div>
                        <h1>🌿 Dashboard Kesehatan</h1>
                        <p>Pantau kesehatan fisik & mental Anda setiap hari</p>
                    </div>
                    <div className="header-actions">
                        <button
                            className="history-btn"
                            onClick={() => setShowHistoryModal(true)}
                        >
                            <History size={18} /> Riwayat
                        </button>
                        <Link to="/health/add" className="add-btn">
                            <Plus size={20} /> Tambah Data
                        </Link>
                    </div>
                </header>

                {/* Stats Grid */}
                <div className="stats-grid">
                    <div className="stat-card health-score">
                        <div className="stat-icon" style={{ background: getScoreColor(healthScore) }}>
                            <Heart />
                        </div>
                        <div className="stat-info">
                            <span className="stat-label">Skor Kesehatan</span>
                            <span className="stat-value" style={{ color: getScoreColor(healthScore) }}>
                                {healthScore}
                            </span>
                        </div>
                        <div className="score-ring" style={{ '--score': healthScore }}>
                            <svg viewBox="0 0 36 36">
                                <path
                                    d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                                    fill="none"
                                    stroke="#e5e7eb"
                                    strokeWidth="3"
                                />
                                <path
                                    d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                                    fill="none"
                                    stroke={getScoreColor(healthScore)}
                                    strokeWidth="3"
                                    strokeDasharray={`${healthScore}, 100`}
                                />
                            </svg>
                        </div>
                    </div>

                    <div className="stat-card bmi-card">
                        <div className="stat-icon" style={{ background: getBMIColor(bmiCategory) }}>
                            <Activity />
                        </div>
                        <div className="stat-info">
                            <span className="stat-label">BMI</span>
                            <span className="stat-value">{bmi.toFixed(1)}</span>
                            <span className="stat-category" style={{ color: getBMIColor(bmiCategory) }}>
                                {bmiCategory}
                            </span>
                        </div>
                    </div>

                    <div className="stat-card mood-card">
                        <div className="stat-icon mood-icon-bg">
                            {getMoodIcon(emotionalState)}
                        </div>
                        <div className="stat-info">
                            <span className="stat-label">Suasana Hati</span>
                            <span className="stat-value mood-value">{getMoodLabel(emotionalState)}</span>
                        </div>
                    </div>

                    <div
                        className="stat-card symptoms-card-clickable"
                        onClick={() => setShowSymptomsModal(true)}
                    >
                        <div className="stat-icon symptoms-icon">
                            <AlertCircle />
                        </div>
                        <div className="stat-info">
                            <span className="stat-label">Gejala Minggu Ini</span>
                            <span className="stat-value">{dashboard?.recent_symptoms?.length || 0}</span>
                        </div>
                        <ChevronRight size={20} className="stat-chevron" />
                    </div>
                </div>

                {/* Charts Section */}
                <div className="charts-section">
                    <div className="chart-card wide">
                        <div className="card-header">
                            <h3><BarChart3 size={20} /> Perkembangan Kesehatan</h3>
                            <div className="period-toggle">
                                <button
                                    className={graphPeriod === 'week' ? 'active' : ''}
                                    onClick={() => setGraphPeriod('week')}
                                >
                                    Mingguan
                                </button>
                                <button
                                    className={graphPeriod === 'month' ? 'active' : ''}
                                    onClick={() => setGraphPeriod('month')}
                                >
                                    Bulanan
                                </button>
                            </div>
                        </div>
                        {graphData.length > 0 ? (
                            <ResponsiveContainer width="100%" height={220}>
                                <AreaChart data={graphData}>
                                    <defs>
                                        <linearGradient id="weightGradient" x1="0" y1="0" x2="0" y2="1">
                                            <stop offset="5%" stopColor="#10b981" stopOpacity={0.3} />
                                            <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
                                        </linearGradient>
                                    </defs>
                                    <XAxis
                                        dataKey="date"
                                        tickFormatter={(d) => d.split('-').slice(1).join('/')}
                                        stroke="#94a3b8"
                                    />
                                    <YAxis domain={['auto', 'auto']} stroke="#94a3b8" />
                                    <Tooltip
                                        contentStyle={{
                                            background: 'rgba(255,255,255,0.95)',
                                            border: '1px solid #10b981',
                                            borderRadius: '12px'
                                        }}
                                    />
                                    <Area
                                        type="monotone"
                                        dataKey="weight"
                                        stroke="#10b981"
                                        strokeWidth={3}
                                        fill="url(#weightGradient)"
                                        dot={{ fill: '#10b981', r: 4 }}
                                        name="Berat (kg)"
                                    />
                                </AreaChart>
                            </ResponsiveContainer>
                        ) : (
                            <div className="empty-chart">
                                <Calendar size={40} />
                                <p>Belum ada data untuk periode ini</p>
                            </div>
                        )}
                    </div>
                </div>

                {/* Dashboard Grid */}
                <div className="dashboard-grid">
                    {/* Reminders Card */}
                    <div className="reminders-card">
                        <div className="card-header">
                            <h3><Bell size={20} /> Pengingat Hari Ini</h3>
                            <div className="reminder-header-actions">
                                <span className="reminder-progress">
                                    {getCompletedReminders()}/{reminders.length} selesai
                                </span>
                                <button className="add-reminder-btn" onClick={openAddReminderModal}>
                                    <Plus size={16} />
                                </button>
                            </div>
                        </div>
                        <div className="reminders-list">
                            {reminders.slice(0, 5).map(reminder => {
                                const IconComponent = getReminderIcon(reminder.type);
                                return (
                                    <div
                                        key={reminder.id}
                                        className={`reminder-item ${!reminder.is_active ? 'done' : ''}`}
                                    >
                                        <div className="reminder-check" onClick={() => toggleReminder(reminder.id)}>
                                            {!reminder.is_active ? <CheckCircle size={18} /> : <div className="check-circle" />}
                                        </div>
                                        <IconComponent size={18} className="reminder-type-icon" />
                                        <span className="reminder-label">{reminder.label}</span>
                                        <span className="reminder-time">
                                            <Clock size={12} /> {reminder.time}
                                        </span>
                                        <div className="reminder-actions">
                                            <button className="reminder-edit-btn" onClick={(e) => { e.stopPropagation(); openEditReminderModal(reminder); }}>
                                                <Edit3 size={14} />
                                            </button>
                                            <button className="reminder-delete-btn" onClick={(e) => { e.stopPropagation(); handleDeleteReminder(reminder.id); }}>
                                                <Trash2 size={14} />
                                            </button>
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                        <button
                            className="see-all-btn"
                            onClick={() => setShowRemindersPanel(true)}
                        >
                            Lihat Semua Pengingat
                        </button>
                    </div>

                    {/* Quick Actions */}
                    <div className="quick-actions">
                        <h3>Menu Cepat</h3>
                        <div className="action-list">
                            <Link to="/symptoms" className="action-item">
                                <Brain className="action-icon" />
                                <span>Log Gejala</span>
                                <ChevronRight size={18} />
                            </Link>
                            <Link to="/recommendations/food" className="action-item">
                                <Utensils className="action-icon food" />
                                <span>Rekomendasi Makanan</span>
                                <ChevronRight size={18} />
                            </Link>
                            <Link to="/recommendations/exercise" className="action-item">
                                <Dumbbell className="action-icon exercise" />
                                <span>Rekomendasi Olahraga</span>
                                <ChevronRight size={18} />
                            </Link>
                            <Link to="/recommendations/emotional" className="action-item">
                                <Heart className="action-icon emotional" />
                                <span>Aktivitas Emosional</span>
                                <ChevronRight size={18} />
                            </Link>
                            <Link to="/family" className="action-item">
                                <Users className="action-icon family" />
                                <span>Keluarga</span>
                                <ChevronRight size={18} />
                            </Link>
                        </div>
                    </div>
                </div>

                {/* Recommendations Preview */}
                {dashboard?.recommendations?.length > 0 && (
                    <div className="recommendations-preview">
                        <h3>💡 Rekomendasi untuk Anda</h3>
                        <div className="rec-list">
                            {dashboard.recommendations.map((rec, idx) => (
                                <div key={idx} className={`rec-card priority-${rec.priority}`}>
                                    <h4>{rec.title}</h4>
                                    <p>{rec.description}</p>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </div>

            {/* Symptoms Modal */}
            {showSymptomsModal && (
                <div className="symptoms-modal-overlay" onClick={() => setShowSymptomsModal(false)}>
                    <div className="symptoms-modal" onClick={(e) => e.stopPropagation()}>
                        <div className="symptoms-modal-header">
                            <h2>
                                <AlertCircle size={24} />
                                Gejala Minggu Ini
                            </h2>
                            <button
                                className="symptoms-modal-close"
                                onClick={() => setShowSymptomsModal(false)}
                            >
                                <X size={24} />
                            </button>
                        </div>
                        <div className="symptoms-modal-footer">
                            <Link
                                to="/symptoms"
                                className="symptoms-modal-btn"
                                onClick={() => setShowSymptomsModal(false)}
                            >
                                <Brain size={18} />
                                Log Gejala Baru
                            </Link>
                        </div>
                        <div className="symptoms-modal-content">
                            {dashboard?.recent_symptoms?.length > 0 ? (
                                <div className="symptoms-list">
                                    {dashboard.recent_symptoms.map((symptom, idx) => {
                                        const getSeverityLevel = (sev) => {
                                            if (sev <= 3) return 1;
                                            if (sev <= 6) return 2;
                                            return 3;
                                        };
                                        const getSeverityLabel = (sev) => {
                                            if (sev <= 3) return 'Ringan';
                                            if (sev <= 6) return 'Sedang';
                                            return 'Berat';
                                        };
                                        const severityLevel = getSeverityLevel(symptom.severity);

                                        return (
                                            <div key={idx} className={`symptom-item severity-${severityLevel}`}>
                                                <div className="symptom-info">
                                                    <span className="symptom-name">{symptom.symptom_name}</span>
                                                    <span className="symptom-type-badge">
                                                        {symptom.symptom_type === 'physical' ? '🩺 Fisik' : '🧠 Mental'}
                                                    </span>
                                                    <span className="symptom-date">
                                                        <Calendar size={14} />
                                                        {new Date(symptom.logged_at).toLocaleDateString('id-ID', {
                                                            weekday: 'long',
                                                            day: 'numeric',
                                                            month: 'long'
                                                        })}
                                                    </span>
                                                </div>
                                                <div className={`symptom-severity severity-badge-${severityLevel}`}>
                                                    {getSeverityLabel(symptom.severity)}
                                                    <span className="severity-number">({symptom.severity}/10)</span>
                                                </div>
                                            </div>
                                        );
                                    })}
                                </div>
                            ) : (
                                <div className="symptoms-empty">
                                    <Heart size={48} />
                                    <p>Tidak ada gejala tercatat minggu ini</p>
                                    <span>Tetap jaga kesehatan Anda!</span>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            )}

            {/* History Modal */}
            {showHistoryModal && (
                <div className="symptoms-modal-overlay" onClick={() => setShowHistoryModal(false)}>
                    <div className="symptoms-modal history-modal" onClick={(e) => e.stopPropagation()}>
                        <div className="symptoms-modal-header">
                            <h2>
                                <History size={24} />
                                Riwayat Kesehatan
                            </h2>
                            <button
                                className="symptoms-modal-close"
                                onClick={() => setShowHistoryModal(false)}
                            >
                                <X size={24} />
                            </button>
                        </div>
                        <div className="symptoms-modal-content">
                            <div className="history-tabs">
                                <h4>📊 Statistik</h4>
                                <div className="history-stats">
                                    <div className="history-stat-item">
                                        <span className="stat-number">{dashboard?.total_records || 0}</span>
                                        <span className="stat-desc">Total Record</span>
                                    </div>
                                    <div className="history-stat-item">
                                        <span className="stat-number">{symptomHistory?.length || dashboard?.recent_symptoms?.length || 0}</span>
                                        <span className="stat-desc">Total Gejala</span>
                                    </div>
                                    <div className="history-stat-item">
                                        <span className="stat-number">{healthScore}</span>
                                        <span className="stat-desc">Skor Saat Ini</span>
                                    </div>
                                </div>
                            </div>
                            {symptomHistory.length > 0 && (
                                <div className="history-section">
                                    <h4>🩺 Riwayat Gejala</h4>
                                    <div className="history-list">
                                        {symptomHistory.slice(0, 10).map((item, idx) => (
                                            <div key={idx} className="history-item">
                                                <span className="history-name">{item.symptom_name}</span>
                                                <span className="history-date">
                                                    {new Date(item.logged_at).toLocaleDateString('id-ID')}
                                                </span>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            )}

            {/* Reminders Panel */}
            {showRemindersPanel && (
                <div className="reminders-panel-overlay" onClick={() => setShowRemindersPanel(false)}>
                    <div className="reminders-panel" onClick={(e) => e.stopPropagation()}>
                        <div className="reminders-panel-header">
                            <h2>
                                <Bell size={24} />
                                Pengingat Kesehatan
                            </h2>
                            <button onClick={() => setShowRemindersPanel(false)}>
                                <X size={24} />
                            </button>
                        </div>
                        <div className="reminders-panel-content">
                            <div className="reminder-categories">
                                <div className="reminder-category">
                                    <h4>💧 Hidrasi</h4>
                                    {reminders.filter(r => r.type === 'water').map(r => (
                                        <div key={r.id} className={`reminder-item ${!r.is_active ? 'done' : ''}`} onClick={() => toggleReminder(r.id)}>
                                            <div className="reminder-check">
                                                {!r.is_active ? <CheckCircle size={16} /> : <div className="check-circle" />}
                                            </div>
                                            <span>{r.label}</span>
                                            <span className="reminder-time">{r.time}</span>
                                        </div>
                                    ))}
                                </div>
                                <div className="reminder-category">
                                    <h4>🍽️ Makan Sehat</h4>
                                    {reminders.filter(r => r.type === 'meal').map(r => (
                                        <div key={r.id} className={`reminder-item ${!r.is_active ? 'done' : ''}`} onClick={() => toggleReminder(r.id)}>
                                            <div className="reminder-check">
                                                {!r.is_active ? <CheckCircle size={16} /> : <div className="check-circle" />}
                                            </div>
                                            <span>{r.label}</span>
                                            <span className="reminder-time">{r.time}</span>
                                        </div>
                                    ))}
                                </div>
                                <div className="reminder-category">
                                    <h4>🏃 Olahraga & Meditasi</h4>
                                    {reminders.filter(r => r.type === 'exercise' || r.type === 'meditation').map(r => (
                                        <div key={r.id} className={`reminder-item ${!r.is_active ? 'done' : ''}`} onClick={() => toggleReminder(r.id)}>
                                            <div className="reminder-check">
                                                {!r.is_active ? <CheckCircle size={16} /> : <div className="check-circle" />}
                                            </div>
                                            <span>{r.label}</span>
                                            <span className="reminder-time">{r.time}</span>
                                        </div>
                                    ))}
                                </div>
                                <div className="reminder-category">
                                    <h4>😴 Istirahat</h4>
                                    {reminders.filter(r => r.type === 'rest').map(r => (
                                        <div key={r.id} className={`reminder-item ${!r.is_active ? 'done' : ''}`} onClick={() => toggleReminder(r.id)}>
                                            <div className="reminder-check">
                                                {!r.is_active ? <CheckCircle size={16} /> : <div className="check-circle" />}
                                            </div>
                                            <span>{r.label}</span>
                                            <span className="reminder-time">{r.time}</span>
                                        </div>
                                    ))}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {/* Reminder Modal */}
            {showReminderModal && (
                <div className="reminder-modal-overlay" onClick={closeReminderModal}>
                    <div className="reminder-modal" onClick={(e) => e.stopPropagation()}>
                        <div className="reminder-modal-header">
                            <h2>
                                {editingReminder ? '✏️ Edit Pengingat' : '➕ Tambah Pengingat'}
                            </h2>
                            <button onClick={closeReminderModal}>
                                <X size={24} />
                            </button>
                        </div>



                        <form onSubmit={handleSaveReminder} className="reminder-modal-form">
                            <div className="form-group">
                                <label>Tipe</label>
                                <select
                                    value={reminderForm.type}
                                    onChange={(e) => setReminderForm({ ...reminderForm, type: e.target.value })}
                                >
                                    <option value="water">💧 Minum Air</option>
                                    <option value="meal">🍽️ Makan</option>
                                    <option value="exercise">🏃 Olahraga</option>
                                    <option value="meditation">🧘 Meditasi</option>
                                    <option value="rest">😴 Istirahat</option>
                                    <option value="custom">⏰ Kustom</option>
                                </select>
                            </div>
                            <div className="form-group">
                                <label>Label</label>
                                <input
                                    type="text"
                                    value={reminderForm.label}
                                    onChange={(e) => setReminderForm({ ...reminderForm, label: e.target.value })}
                                    placeholder="Contoh: Minum Air Pagi"
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label>Waktu</label>
                                <input
                                    type="time"
                                    value={reminderForm.time}
                                    onChange={(e) => setReminderForm({ ...reminderForm, time: e.target.value })}
                                    required
                                />
                            </div>
                            <div className="reminder-modal-actions">
                                <button type="button" className="cancel-btn" onClick={closeReminderModal}>
                                    Batal
                                </button>
                                <button type="submit" className="save-btn" disabled={reminderLoading}>
                                    {reminderLoading ? 'Menyimpan...' : (editingReminder ? 'Update' : 'Simpan')}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Dashboard;
