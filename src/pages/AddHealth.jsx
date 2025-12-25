import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { healthAPI } from '../services/api';
import {
    ArrowLeft, Scale, Ruler, Activity, Smile,
    FileText, Loader2, Check
} from 'lucide-react';
import './AddHealth.css';

const AddHealth = () => {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [success, setSuccess] = useState(false);

    const [formData, setFormData] = useState({
        heightCm: '',
        weightKg: '',
        activityLevel: 'sedentary',
        emotionalState: 'neutral',
        notes: '',
    });

    const activityLevels = [
        { value: 'sedentary', label: 'Tidak Aktif' },
        { value: 'light', label: 'Ringan' },
        { value: 'moderate', label: 'Sedang' },
        { value: 'active', label: 'Aktif' },
    ];

    const emotionalStates = [
        { value: 'happy', label: '😊 Bahagia' },
        { value: 'neutral', label: '😐 Biasa' },
        { value: 'stressed', label: '😰 Stres' },
        { value: 'anxious', label: '😟 Cemas' },
        { value: 'sad', label: '😢 Sedih' },
    ];

    const handleChange = (field, value) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);

        try {
            await healthAPI.create({
                height_cm: parseFloat(formData.heightCm),
                weight_kg: parseFloat(formData.weightKg),
                activity_level: formData.activityLevel,
                emotional_state: formData.emotionalState,
                notes: formData.notes,
            });

            setSuccess(true);
            setTimeout(() => {
                navigate('/dashboard');
            }, 1500);
        } catch (err) {
            console.error('Failed to save health data:', err);
        } finally {
            setLoading(false);
        }
    };

    const bmi = formData.heightCm && formData.weightKg
        ? (parseFloat(formData.weightKg) / Math.pow(parseFloat(formData.heightCm) / 100, 2)).toFixed(1)
        : null;

    return (
        <div className="add-health-page">
            {/* Animated Health Background */}
            <div className="health-bg-elements">
                {/* Floating Health Icons */}
                <div className="health-particle">❤️</div>
                <div className="health-particle">💪</div>
                <div className="health-particle">🏃</div>
                <div className="health-particle">🍎</div>
                <div className="health-particle">💧</div>
                <div className="health-particle">🧘</div>
                {/* Pulse Rings */}
                <div className="health-pulse"></div>
                <div className="health-pulse"></div>
                <div className="health-pulse"></div>
            </div>

            <header className="page-header">
                <Link to="/dashboard" className="back-btn">
                    <ArrowLeft size={20} />
                </Link>
                <h1>Tambah Data Kesehatan</h1>
            </header>

            <div className="add-health-container">
                <form onSubmit={handleSubmit}>
                    <div className="form-section">
                        <h3><Scale size={18} /> Data Fisik</h3>

                        <div className="form-row">
                            <div className="form-group">
                                <label>Tinggi Badan (cm)</label>
                                <input
                                    type="number"
                                    placeholder="170"
                                    value={formData.heightCm}
                                    onChange={(e) => handleChange('heightCm', e.target.value)}
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label>Berat Badan (kg)</label>
                                <input
                                    type="number"
                                    placeholder="65"
                                    value={formData.weightKg}
                                    onChange={(e) => handleChange('weightKg', e.target.value)}
                                    required
                                />
                            </div>
                        </div>

                        {bmi && (
                            <div className="bmi-preview">
                                <span>BMI:</span>
                                <strong>{bmi}</strong>
                            </div>
                        )}
                    </div>

                    <div className="form-section">
                        <h3><Activity size={18} /> Tingkat Aktivitas</h3>
                        <div className="option-row">
                            {activityLevels.map((level) => (
                                <button
                                    key={level.value}
                                    type="button"
                                    className={`option-btn ${formData.activityLevel === level.value ? 'selected' : ''}`}
                                    onClick={() => handleChange('activityLevel', level.value)}
                                >
                                    {level.label}
                                </button>
                            ))}
                        </div>
                    </div>

                    <div className="form-section">
                        <h3><Smile size={18} /> Kondisi Emosional</h3>
                        <div className="emotion-row">
                            {emotionalStates.map((state) => (
                                <button
                                    key={state.value}
                                    type="button"
                                    className={`emotion-btn ${formData.emotionalState === state.value ? 'selected' : ''}`}
                                    onClick={() => handleChange('emotionalState', state.value)}
                                >
                                    {state.label}
                                </button>
                            ))}
                        </div>
                    </div>

                    <div className="form-section">
                        <h3><FileText size={18} /> Catatan (Opsional)</h3>
                        <textarea
                            placeholder="Catatan tambahan tentang kondisi Anda hari ini..."
                            value={formData.notes}
                            onChange={(e) => handleChange('notes', e.target.value)}
                            rows={3}
                        />
                    </div>

                    <button type="submit" className="submit-btn" disabled={loading || success}>
                        {loading ? (
                            <Loader2 className="spinner" />
                        ) : success ? (
                            <>
                                <Check size={20} /> Tersimpan!
                            </>
                        ) : (
                            'Simpan Data'
                        )}
                    </button>
                </form>
            </div>
        </div>
    );
};

export default AddHealth;
