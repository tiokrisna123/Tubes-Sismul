import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { symptomsAPI } from '../services/api';
import {
    ArrowLeft, Check, AlertCircle, Activity, Brain,
    Loader2, CheckCircle, Plus
} from 'lucide-react';
import './Symptoms.css';

const Symptoms = () => {
    const [symptoms, setSymptoms] = useState({ physical: [], mental: [] });
    const [selected, setSelected] = useState([]);
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [success, setSuccess] = useState(false);
    const [activeTab, setActiveTab] = useState('physical');
    const [customSymptom, setCustomSymptom] = useState('');

    useEffect(() => {
        loadSymptoms();
    }, []);

    const loadSymptoms = async () => {
        try {
            const res = await symptomsAPI.getList();
            setSymptoms(res.data.data);
        } catch (err) {
            console.error('Failed to load symptoms:', err);
        } finally {
            setLoading(false);
        }
    };

    const toggleSymptom = (symptom, type) => {
        const exists = selected.find(s => s.symptom_name === symptom.symptom_name);
        if (exists) {
            setSelected(selected.filter(s => s.symptom_name !== symptom.symptom_name));
        } else {
            setSelected([...selected, {
                symptom_type: type,
                symptom_name: symptom.symptom_name,
                severity: 5,
                notes: ''
            }]);
        }
    };

    const updateSeverity = (name, severity) => {
        setSelected(selected.map(s =>
            s.symptom_name === name ? { ...s, severity } : s
        ));
    };

    const addCustomSymptom = () => {
        if (!customSymptom.trim()) return;
        setSelected([...selected, {
            symptom_type: activeTab,
            symptom_name: customSymptom.trim(),
            severity: 5,
            notes: ''
        }]);
        setCustomSymptom('');
    };

    const handleSubmit = async () => {
        if (selected.length === 0) return;
        setSubmitting(true);
        try {
            await symptomsAPI.logBatch(selected);
            setSuccess(true);
            setTimeout(() => {
                setSuccess(false);
                setSelected([]);
            }, 2000);
        } catch (err) {
            console.error('Failed to log symptoms:', err);
        } finally {
            setSubmitting(false);
        }
    };

    if (loading) {
        return (
            <div className="symptoms-loading">
                <Loader2 className="spinner" />
                <p>Memuat daftar gejala...</p>
            </div>
        );
    }

    return (
        <div className="symptoms-page">
            <header className="page-header">
                <Link to="/dashboard" className="back-btn">
                    <ArrowLeft size={20} />
                </Link>
                <h1>Log Gejala</h1>
            </header>

            <div className="symptoms-container">
                <div className="tabs">
                    <button
                        className={`tab ${activeTab === 'physical' ? 'active' : ''}`}
                        onClick={() => setActiveTab('physical')}
                    >
                        <Activity size={18} /> Fisik
                    </button>
                    <button
                        className={`tab ${activeTab === 'mental' ? 'active' : ''}`}
                        onClick={() => setActiveTab('mental')}
                    >
                        <Brain size={18} /> Mental
                    </button>
                </div>

                <div className="symptom-grid">
                    {symptoms[activeTab]?.map((symptom) => {
                        const isSelected = selected.find(s => s.symptom_name === symptom.symptom_name);
                        return (
                            <button
                                key={symptom.id}
                                className={`symptom-chip ${isSelected ? 'selected' : ''}`}
                                onClick={() => toggleSymptom(symptom, activeTab)}
                            >
                                {isSelected && <CheckCircle size={16} />}
                                {symptom.symptom_name}
                            </button>
                        );
                    })}
                </div>

                <div className="custom-symptom">
                    <input
                        type="text"
                        placeholder="Tambah gejala lain..."
                        value={customSymptom}
                        onChange={(e) => setCustomSymptom(e.target.value)}
                        onKeyPress={(e) => e.key === 'Enter' && addCustomSymptom()}
                    />
                    <button onClick={addCustomSymptom}>
                        <Plus size={18} />
                    </button>
                </div>

                {selected.length > 0 && (
                    <div className="selected-symptoms">
                        <h3>Gejala Dipilih ({selected.length})</h3>
                        {selected.map((s) => (
                            <div key={s.symptom_name} className="selected-item">
                                <span className="selected-name">{s.symptom_name}</span>
                                <div className="severity-slider">
                                    <span>Tingkat: {s.severity}</span>
                                    <input
                                        type="range"
                                        min="1"
                                        max="10"
                                        value={s.severity}
                                        onChange={(e) => updateSeverity(s.symptom_name, parseInt(e.target.value))}
                                    />
                                </div>
                            </div>
                        ))}
                    </div>
                )}

                <button
                    className="submit-btn"
                    onClick={handleSubmit}
                    disabled={selected.length === 0 || submitting}
                >
                    {submitting ? (
                        <Loader2 className="spinner" />
                    ) : success ? (
                        <>
                            <Check size={20} /> Tersimpan!
                        </>
                    ) : (
                        <>
                            <AlertCircle size={20} /> Simpan Gejala
                        </>
                    )}
                </button>
            </div>
        </div>
    );
};

export default Symptoms;
