import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { authAPI, healthAPI } from '../services/api';
import {
    Heart, ArrowRight, ArrowLeft, Scale, Ruler, Activity,
    Smile, Calendar, Loader2, CheckCircle
} from 'lucide-react';
import './Onboarding.css';

const Onboarding = () => {
    const [step, setStep] = useState(1);
    const [loading, setLoading] = useState(false);
    const { updateUser } = useAuth();
    const navigate = useNavigate();

    const [formData, setFormData] = useState({
        birthDate: '',
        heightCm: '',
        weightKg: '',
        activityLevel: 'sedentary',
        emotionalState: 'neutral',
        dailySchedule: '',
    });

    const activityLevels = [
        { value: 'sedentary', label: 'Tidak Aktif', desc: 'Kerja kantoran, jarang olahraga' },
        { value: 'light', label: 'Ringan', desc: 'Olahraga 1-2x/minggu' },
        { value: 'moderate', label: 'Sedang', desc: 'Olahraga 3-4x/minggu' },
        { value: 'active', label: 'Aktif', desc: 'Olahraga 5+x/minggu' },
    ];

    const emotionalStates = [
        { value: 'happy', label: '😊 Bahagia', color: '#48bb78' },
        { value: 'neutral', label: '😐 Biasa', color: '#a0aec0' },
        { value: 'stressed', label: '😰 Stres', color: '#ed8936' },
        { value: 'anxious', label: '😟 Cemas', color: '#f56565' },
        { value: 'sad', label: '😢 Sedih', color: '#667eea' },
    ];

    const handleChange = (field, value) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    };

    const handleSubmit = async () => {
        setLoading(true);
        try {
            // Update profile
            const profileRes = await authAPI.updateProfile({
                birth_date: formData.birthDate ? new Date(formData.birthDate).toISOString() : null,
                height_cm: parseFloat(formData.heightCm),
                weight_kg: parseFloat(formData.weightKg),
                activity_level: formData.activityLevel,
            });
            updateUser(profileRes.data.data);

            // Create initial health data
            await healthAPI.create({
                height_cm: parseFloat(formData.heightCm),
                weight_kg: parseFloat(formData.weightKg),
                activity_level: formData.activityLevel,
                emotional_state: formData.emotionalState,
                daily_schedule: formData.dailySchedule,
            });

            navigate('/dashboard');
        } catch (err) {
            console.error('Onboarding error:', err);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="onboarding-container">
            <div className="onboarding-card">
                <div className="onboarding-header">
                    <div className="step-indicator">
                        {[1, 2, 3].map((s) => (
                            <div key={s} className={`step-dot ${step >= s ? 'active' : ''}`}>
                                {step > s ? <CheckCircle size={16} /> : s}
                            </div>
                        ))}
                    </div>
                </div>

                {step === 1 && (
                    <div className="onboarding-step">
                        <div className="step-icon">
                            <Scale />
                        </div>
                        <h2>Data Fisik Dasar</h2>
                        <p>Masukkan informasi dasar untuk menghitung BMI dan rekomendasi kesehatan</p>

                        <div className="form-group">
                            <label><Calendar size={18} /> Tanggal Lahir</label>
                            <input
                                type="date"
                                value={formData.birthDate}
                                onChange={(e) => handleChange('birthDate', e.target.value)}
                            />
                        </div>

                        <div className="form-row">
                            <div className="form-group">
                                <label><Ruler size={18} /> Tinggi (cm)</label>
                                <input
                                    type="number"
                                    placeholder="Contoh: 170"
                                    value={formData.heightCm}
                                    onChange={(e) => handleChange('heightCm', e.target.value)}
                                />
                            </div>
                            <div className="form-group">
                                <label><Scale size={18} /> Berat (kg)</label>
                                <input
                                    type="number"
                                    placeholder="Contoh: 65"
                                    value={formData.weightKg}
                                    onChange={(e) => handleChange('weightKg', e.target.value)}
                                />
                            </div>
                        </div>
                    </div>
                )}

                {step === 2 && (
                    <div className="onboarding-step">
                        <div className="step-icon">
                            <Activity />
                        </div>
                        <h2>Tingkat Aktivitas</h2>
                        <p>Pilih yang paling sesuai dengan kegiatan harian Anda</p>

                        <div className="option-grid">
                            {activityLevels.map((level) => (
                                <button
                                    key={level.value}
                                    className={`option-card ${formData.activityLevel === level.value ? 'selected' : ''}`}
                                    onClick={() => handleChange('activityLevel', level.value)}
                                >
                                    <span className="option-label">{level.label}</span>
                                    <span className="option-desc">{level.desc}</span>
                                </button>
                            ))}
                        </div>
                    </div>
                )}

                {step === 3 && (
                    <div className="onboarding-step">
                        <div className="step-icon">
                            <Smile />
                        </div>
                        <h2>Kondisi Emosional</h2>
                        <p>Bagaimana perasaan Anda saat ini?</p>

                        <div className="emotion-grid">
                            {emotionalStates.map((emotion) => (
                                <button
                                    key={emotion.value}
                                    className={`emotion-card ${formData.emotionalState === emotion.value ? 'selected' : ''}`}
                                    onClick={() => handleChange('emotionalState', emotion.value)}
                                    style={{ '--emotion-color': emotion.color }}
                                >
                                    {emotion.label}
                                </button>
                            ))}
                        </div>

                        <div className="form-group">
                            <label>Ceritakan rutinitas harian Anda (opsional)</label>
                            <textarea
                                placeholder="Contoh: Kerja jam 9-17, tidur jam 23:00..."
                                value={formData.dailySchedule}
                                onChange={(e) => handleChange('dailySchedule', e.target.value)}
                                rows={3}
                            />
                        </div>
                    </div>
                )}

                <div className="onboarding-actions">
                    {step > 1 && (
                        <button className="btn-secondary" onClick={() => setStep(step - 1)}>
                            <ArrowLeft size={18} /> Kembali
                        </button>
                    )}

                    {step < 3 ? (
                        <button
                            className="btn-primary"
                            onClick={() => setStep(step + 1)}
                            disabled={step === 1 && (!formData.heightCm || !formData.weightKg)}
                        >
                            Lanjut <ArrowRight size={18} />
                        </button>
                    ) : (
                        <button className="btn-primary" onClick={handleSubmit} disabled={loading}>
                            {loading ? <Loader2 className="spinner" /> : 'Selesai'}
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
};

export default Onboarding;
