import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import { useToast } from '../components/Toast';
import { authAPI } from '../services/api';
import {
    User, Mail, Calendar, Ruler, Scale, Activity,
    ArrowLeft, Save, Moon, Sun, LogOut, Camera
} from 'lucide-react';
import './Profile.css';

const Profile = () => {
    const navigate = useNavigate();
    const { user, logout, updateUser } = useAuth();
    const { isDark, toggleTheme } = useTheme();
    const toast = useToast();

    const [loading, setLoading] = useState(false);
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        birthDate: '',
        heightCm: '',
        weightKg: '',
        activityLevel: 'moderate'
    });

    useEffect(() => {
        if (user) {
            setFormData({
                name: user.name || '',
                email: user.email || '',
                birthDate: user.birth_date ? user.birth_date.split('T')[0] : '',
                heightCm: user.height_cm || '',
                weightKg: user.weight_kg || '',
                activityLevel: user.activity_level || 'moderate'
            });
        }
    }, [user]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);

        try {
            const response = await authAPI.updateProfile({
                name: formData.name,
                birth_date: formData.birthDate,
                height_cm: parseFloat(formData.heightCm) || 0,
                weight_kg: parseFloat(formData.weightKg) || 0,
                activity_level: formData.activityLevel
            });

            if (response.data.success) {
                updateUser(response.data.data);
                toast.success('Profil berhasil diperbarui!');
            }
        } catch (error) {
            toast.error('Gagal memperbarui profil');
            console.error('Update profile error:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    const handleExportData = () => {
        // Create report content
        const reportContent = `
LAPORAN KESEHATAN
=================
Nama: ${user?.name || '-'}
Email: ${user?.email || '-'}
Tanggal: ${new Date().toLocaleDateString('id-ID')}

DATA FISIK
----------
Tinggi: ${formData.heightCm} cm
Berat: ${formData.weightKg} kg
BMI: ${formData.heightCm && formData.weightKg ?
                (formData.weightKg / Math.pow(formData.heightCm / 100, 2)).toFixed(1) : '-'}
Aktivitas: ${formData.activityLevel}

---
Digenerate oleh Health Tracker
    `;

        // Create and download file
        const blob = new Blob([reportContent], { type: 'text/plain;charset=utf-8' });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `health-report-${new Date().toISOString().split('T')[0]}.txt`;
        link.click();
        URL.revokeObjectURL(url);

        toast.success('Laporan berhasil diunduh!');
    };

    return (
        <div className="profile-page">
            <div className="profile-container">
                <header className="profile-header">
                    <button className="back-btn" onClick={() => navigate('/dashboard')}>
                        <ArrowLeft size={20} />
                        Kembali
                    </button>
                    <h1>Pengaturan Profil</h1>
                </header>

                <div className="profile-content">
                    {/* Avatar Section */}
                    <div className="avatar-section">
                        <div className="avatar">
                            <User size={48} />
                        </div>
                        <div className="avatar-info">
                            <h2>{user?.name || 'User'}</h2>
                            <p>{user?.email}</p>
                        </div>
                    </div>

                    {/* Settings Form */}
                    <form onSubmit={handleSubmit} className="profile-form">
                        <div className="form-section">
                            <h3>Informasi Dasar</h3>

                            <div className="form-group">
                                <label><User size={16} /> Nama Lengkap</label>
                                <input
                                    type="text"
                                    name="name"
                                    value={formData.name}
                                    onChange={handleChange}
                                    placeholder="Nama lengkap Anda"
                                />
                            </div>

                            <div className="form-group">
                                <label><Mail size={16} /> Email</label>
                                <input
                                    type="email"
                                    name="email"
                                    value={formData.email}
                                    disabled
                                    className="disabled"
                                />
                                <small>Email tidak dapat diubah</small>
                            </div>

                            <div className="form-group">
                                <label><Calendar size={16} /> Tanggal Lahir</label>
                                <input
                                    type="date"
                                    name="birthDate"
                                    value={formData.birthDate}
                                    onChange={handleChange}
                                />
                            </div>
                        </div>

                        <div className="form-section">
                            <h3>Data Fisik</h3>

                            <div className="form-row">
                                <div className="form-group">
                                    <label><Ruler size={16} /> Tinggi (cm)</label>
                                    <input
                                        type="number"
                                        name="heightCm"
                                        value={formData.heightCm}
                                        onChange={handleChange}
                                        placeholder="170"
                                    />
                                </div>

                                <div className="form-group">
                                    <label><Scale size={16} /> Berat (kg)</label>
                                    <input
                                        type="number"
                                        name="weightKg"
                                        value={formData.weightKg}
                                        onChange={handleChange}
                                        placeholder="65"
                                    />
                                </div>
                            </div>

                            <div className="form-group">
                                <label><Activity size={16} /> Level Aktivitas</label>
                                <select
                                    name="activityLevel"
                                    value={formData.activityLevel}
                                    onChange={handleChange}
                                >
                                    <option value="sedentary">Tidak Aktif (Jarang olahraga)</option>
                                    <option value="light">Ringan (1-2x/minggu)</option>
                                    <option value="moderate">Sedang (3-4x/minggu)</option>
                                    <option value="active">Aktif (5-6x/minggu)</option>
                                    <option value="very_active">Sangat Aktif (Setiap hari)</option>
                                </select>
                            </div>
                        </div>

                        <button type="submit" className="save-btn" disabled={loading}>
                            <Save size={18} />
                            {loading ? 'Menyimpan...' : 'Simpan Perubahan'}
                        </button>
                    </form>

                    {/* Preferences Section */}
                    <div className="preferences-section">
                        <h3>Preferensi</h3>

                        <div className="preference-item">
                            <div className="preference-info">
                                {isDark ? <Moon size={20} /> : <Sun size={20} />}
                                <div>
                                    <span className="preference-title">Mode Gelap</span>
                                    <span className="preference-desc">Menggunakan tema gelap untuk tampilan</span>
                                </div>
                            </div>
                            <button
                                className={`toggle-btn ${isDark ? 'active' : ''}`}
                                onClick={toggleTheme}
                                aria-label="Toggle dark mode"
                            >
                                <span className="toggle-slider"></span>
                            </button>
                        </div>
                    </div>

                    {/* Actions Section */}
                    <div className="actions-section">
                        <button className="export-btn" onClick={handleExportData}>
                            📄 Export Laporan Kesehatan
                        </button>

                        <button className="logout-btn" onClick={handleLogout}>
                            <LogOut size={18} />
                            Keluar dari Akun
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Profile;
