import { useState, useEffect } from 'react';
import { waterAPI } from '../services/api';
import { Droplets, Plus, Minus, Target, Loader2 } from 'lucide-react';
import './WaterTracker.css';

const WaterTracker = () => {
    const [water, setWater] = useState(null);
    const [loading, setLoading] = useState(true);
    const [updating, setUpdating] = useState(false);

    useEffect(() => {
        fetchWater();
    }, []);

    const fetchWater = async () => {
        try {
            const res = await waterAPI.get();
            setWater(res.data);
        } catch (error) {
            console.error('Failed to fetch water:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleAdd = async () => {
        if (updating) return;
        setUpdating(true);
        try {
            const res = await waterAPI.addGlass();
            setWater(res.data);
        } catch (error) {
            console.error('Failed to add glass:', error);
        } finally {
            setUpdating(false);
        }
    };

    const handleRemove = async () => {
        if (updating || !water || water.glasses <= 0) return;
        setUpdating(true);
        try {
            const res = await waterAPI.removeGlass();
            setWater(res.data);
        } catch (error) {
            console.error('Failed to remove glass:', error);
        } finally {
            setUpdating(false);
        }
    };

    if (loading) {
        return (
            <div className="water-tracker loading">
                <Loader2 className="spinner" size={24} />
            </div>
        );
    }

    const glasses = water?.glasses || 0;
    const goal = water?.goal || 8;
    const percentage = water?.percentage || 0;

    return (
        <div className="water-tracker">
            <div className="water-header">
                <div className="water-icon">
                    <Droplets size={24} />
                </div>
                <div className="water-info">
                    <h3>💧 Minum Air</h3>
                    <p>{glasses} dari {goal} gelas</p>
                </div>
            </div>

            <div className="water-progress">
                <div className="water-bar">
                    <div
                        className="water-fill"
                        style={{ width: `${percentage}%` }}
                    />
                </div>
                <span className="water-percentage">{Math.round(percentage)}%</span>
            </div>

            <div className="water-glasses">
                {Array.from({ length: goal }).map((_, index) => (
                    <div
                        key={index}
                        className={`glass ${index < glasses ? 'filled' : ''}`}
                    >
                        <Droplets size={16} />
                    </div>
                ))}
            </div>

            <div className="water-actions">
                <button
                    className="water-btn remove"
                    onClick={handleRemove}
                    disabled={updating || glasses <= 0}
                >
                    <Minus size={20} />
                </button>
                <span className="water-count">{glasses}</span>
                <button
                    className="water-btn add"
                    onClick={handleAdd}
                    disabled={updating}
                >
                    <Plus size={20} />
                </button>
            </div>

            {percentage >= 100 && (
                <div className="water-complete">
                    🎉 Target tercapai!
                </div>
            )}
        </div>
    );
};

export default WaterTracker;
