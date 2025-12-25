import { useState, useEffect } from 'react';
import { Link, useParams } from 'react-router-dom';
import { familyAPI } from '../services/api';
import {
    ArrowLeft, Heart, Activity, AlertCircle, Loader2,
    Calendar
} from 'lucide-react';
import './FamilyHealth.css';

const FamilyHealth = () => {
    const { id } = useParams();
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        loadHealth();
    }, [id]);

    const loadHealth = async () => {
        try {
            const res = await familyAPI.getMemberHealth(id);
            setData(res.data.data);
        } catch (err) {
            setError(err.response?.data?.error || 'Gagal memuat data kesehatan');
        } finally {
            setLoading(false);
        }
    };

    const getBMIColor = (category) => {
        switch (category) {
            case 'Underweight': return '#ed8936';
            case 'Normal': return '#48bb78';
            case 'Overweight': return '#ed8936';
            case 'Obese': return '#f56565';
            default: return '#a0aec0';
        }
    };

    if (loading) {
        return (
            <div className="fh-loading">
                <Loader2 className="spinner" />
                <p>Memuat data kesehatan...</p>
            </div>
        );
    }

    if (error) {
        return (
            <div className="fh-error">
                <AlertCircle size={48} />
                <h3>Tidak Dapat Mengakses</h3>
                <p>{error}</p>
                <Link to="/family" className="back-link">Kembali</Link>
            </div>
        );
    }

    const health = data?.latest_health;
    const bmiCategory = data?.bmi_category || 'Unknown';

    return (
        <div className="family-health-page">
            <header className="page-header">
                <Link to="/family" className="back-btn">
                    <ArrowLeft size={20} />
                </Link>
                <h1>Kesehatan {data?.member_name}</h1>
            </header>

            <div className="fh-container">
                <div className="member-header">
                    <div className="member-avatar large">
                        {data?.member_name?.charAt(0).toUpperCase()}
                    </div>
                    <div>
                        <h2>{data?.member_name}</h2>
                        <span className="relation-tag">{data?.relationship}</span>
                    </div>
                </div>

                {health ? (
                    <>
                        <div className="health-stats">
                            <div className="health-stat">
                                <Activity size={20} />
                                <div>
                                    <span className="label">BMI</span>
                                    <span className="value">{health.bmi?.toFixed(1)}</span>
                                    <span className="category" style={{ color: getBMIColor(bmiCategory) }}>
                                        {bmiCategory}
                                    </span>
                                </div>
                            </div>

                            <div className="health-stat">
                                <Heart size={20} />
                                <div>
                                    <span className="label">Berat</span>
                                    <span className="value">{health.weight_kg} kg</span>
                                </div>
                            </div>

                            <div className="health-stat">
                                <Calendar size={20} />
                                <div>
                                    <span className="label">Emosi</span>
                                    <span className="value emotion">{health.emotional_state || '-'}</span>
                                </div>
                            </div>
                        </div>

                        <div className="last-update">
                            Terakhir update: {new Date(health.record_date).toLocaleDateString('id-ID')}
                        </div>
                    </>
                ) : (
                    <div className="no-data">
                        <Heart size={40} />
                        <p>Belum ada data kesehatan</p>
                    </div>
                )}

                {data?.recent_symptoms?.length > 0 && (
                    <div className="symptoms-section">
                        <h3><AlertCircle size={18} /> Gejala Terkini</h3>
                        <div className="symptom-list">
                            {data.recent_symptoms.map((s) => (
                                <div key={s.id} className="symptom-tag">
                                    <span className={`type-badge ${s.symptom_type}`}>
                                        {s.symptom_type === 'physical' ? '🏥' : '🧠'}
                                    </span>
                                    {s.symptom_name}
                                    <span className="severity">Lvl {s.severity}</span>
                                </div>
                            ))}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default FamilyHealth;
