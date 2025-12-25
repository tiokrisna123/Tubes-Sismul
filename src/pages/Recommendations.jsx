import { useState, useEffect } from 'react';
import { Link, useParams } from 'react-router-dom';
import { recommendationsAPI } from '../services/api';
import {
    ArrowLeft, Utensils, Dumbbell, Heart, Loader2,
    CheckCircle, XCircle, Clock, Flame, Calendar, Coffee, Moon, Sun,
    Apple, Lightbulb, Wallet
} from 'lucide-react';
import './Recommendations.css';

const Recommendations = () => {
    const { type } = useParams();
    const [recommendations, setRecommendations] = useState([]);
    const [dailyMenu, setDailyMenu] = useState(null);
    const [loading, setLoading] = useState(true);
    const [activeTab, setActiveTab] = useState('recommendations');
    const [menuError, setMenuError] = useState(null);

    useEffect(() => {
        loadRecommendations();
    }, [type]);

    const loadRecommendations = async () => {
        setLoading(true);
        try {
            let res;
            switch (type) {
                case 'food':
                    res = await recommendationsAPI.getFood();
                    // Also load daily menu for food type
                    try {
                        console.log('Loading daily menu...');
                        const menuRes = await recommendationsAPI.getDailyMenu();
                        console.log('Daily menu response:', menuRes.data);
                        setDailyMenu(menuRes.data.data);
                        setMenuError(null);
                    } catch (e) {
                        console.error('Failed to load daily menu:', e);
                        const errorMsg = e.response?.data?.error || e.message || 'Gagal memuat menu';
                        setMenuError(errorMsg);
                        setDailyMenu(null);
                    }
                    break;
                case 'exercise':
                    res = await recommendationsAPI.getExercise();
                    break;
                case 'emotional':
                    res = await recommendationsAPI.getEmotional();
                    break;
                default:
                    res = await recommendationsAPI.getFood();
            }
            setRecommendations(res.data.data || []);
        } catch (err) {
            console.error('Failed to load recommendations:', err);
        } finally {
            setLoading(false);
        }
    };

    const getIcon = () => {
        switch (type) {
            case 'food': return <Utensils />;
            case 'exercise': return <Dumbbell />;
            case 'emotional': return <Heart />;
            default: return <Utensils />;
        }
    };

    const getTitle = () => {
        switch (type) {
            case 'food': return 'Rekomendasi Makanan';
            case 'exercise': return 'Rekomendasi Olahraga';
            case 'emotional': return 'Aktivitas Emosional';
            default: return 'Rekomendasi';
        }
    };

    const renderMealCard = (meal, icon) => (
        <div className="meal-card">
            <div className="meal-header">
                {icon}
                <div>
                    <h4>{meal.title}</h4>
                    <div className="meal-meta">
                        <span className="meal-calories">{meal.calories}</span>
                        {meal.estimated_cost && (
                            <span className="meal-cost"><Wallet size={12} /> {meal.estimated_cost}</span>
                        )}
                    </div>
                </div>
            </div>
            <p className="meal-description">{meal.description}</p>
            <div className="meal-foods">
                {meal.foods?.map((food, i) => (
                    <span key={i} className="food-item">{food}</span>
                ))}
            </div>
            {meal.ingredients && meal.ingredients.length > 0 && (
                <div className="meal-ingredients">
                    <strong>🛒 Bahan-bahan:</strong>
                    <ul>
                        {meal.ingredients.map((ing, i) => (
                            <li key={i}>{ing}</li>
                        ))}
                    </ul>
                </div>
            )}
            {meal.recipe && (
                <div className="meal-recipe">
                    <strong>👨‍🍳 Cara Membuat:</strong>
                    <p>{meal.recipe}</p>
                </div>
            )}
        </div>
    );

    if (loading) {
        return (
            <div className="rec-loading">
                <Loader2 className="spinner" />
                <p>Memuat rekomendasi...</p>
            </div>
        );
    }

    return (
        <div className="recommendations-page">
            <header className="page-header">
                <Link to="/dashboard" className="back-btn">
                    <ArrowLeft size={20} />
                </Link>
                <div className="header-icon">{getIcon()}</div>
                <h1>{getTitle()}</h1>
            </header>

            <div className="rec-container">
                {/* Tab navigation for food type */}
                {type === 'food' && (
                    <div className="rec-tabs">
                        <button
                            className={`rec-tab ${activeTab === 'recommendations' ? 'active' : ''}`}
                            onClick={() => setActiveTab('recommendations')}
                        >
                            <Utensils size={18} />
                            Rekomendasi
                        </button>
                        <button
                            className={`rec-tab ${activeTab === 'dailyMenu' ? 'active' : ''}`}
                            onClick={() => setActiveTab('dailyMenu')}
                        >
                            <Calendar size={18} />
                            Menu Harian
                        </button>
                    </div>
                )}

                {/* Daily Menu Section */}
                {type === 'food' && activeTab === 'dailyMenu' && (
                    <div className="daily-menu-section">
                        {dailyMenu ? (
                            <>
                                <div className="menu-header-card">
                                    <div className="menu-date">
                                        <Calendar size={24} />
                                        <div>
                                            <h3>Menu {dailyMenu.date}</h3>
                                            <div className="menu-totals">
                                                <span className="total-cal">🔥 {dailyMenu.total_calories}</span>
                                                {dailyMenu.total_estimated_cost && (
                                                    <span className="total-cost"><Wallet size={14} /> {dailyMenu.total_estimated_cost}</span>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                    <div className="health-tip">
                                        <Lightbulb size={18} />
                                        <p>{dailyMenu.health_tip}</p>
                                    </div>
                                </div>

                                <div className="meals-grid">
                                    {dailyMenu.breakfast && renderMealCard(dailyMenu.breakfast, <Coffee size={24} className="meal-icon breakfast" />)}
                                    {dailyMenu.lunch && renderMealCard(dailyMenu.lunch, <Sun size={24} className="meal-icon lunch" />)}
                                    {dailyMenu.dinner && renderMealCard(dailyMenu.dinner, <Moon size={24} className="meal-icon dinner" />)}
                                </div>

                                {/* Alternative Breakfast Options */}
                                {dailyMenu.breakfast_alt && dailyMenu.breakfast_alt.length > 0 && (
                                    <div className="alt-meals-section">
                                        <h3>☀️ Alternatif Sarapan Lainnya</h3>
                                        <div className="alt-meals-grid">
                                            {dailyMenu.breakfast_alt.map((meal, idx) => (
                                                <div key={idx} className="alt-meal-card">
                                                    <h4>{meal.title}</h4>
                                                    <p className="alt-meal-desc">{meal.description}</p>
                                                    <div className="alt-meal-foods">
                                                        {meal.foods?.slice(0, 3).map((food, i) => (
                                                            <span key={i}>{food}</span>
                                                        ))}
                                                    </div>
                                                    <div className="alt-meal-meta">
                                                        <span className="alt-cal">{meal.calories}</span>
                                                        <span className="alt-cost">{meal.estimated_cost}</span>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}

                                {/* Alternative Lunch Options */}
                                {dailyMenu.lunch_alt && dailyMenu.lunch_alt.length > 0 && (
                                    <div className="alt-meals-section">
                                        <h3>🌤️ Alternatif Makan Siang Lainnya</h3>
                                        <div className="alt-meals-grid">
                                            {dailyMenu.lunch_alt.map((meal, idx) => (
                                                <div key={idx} className="alt-meal-card">
                                                    <h4>{meal.title}</h4>
                                                    <p className="alt-meal-desc">{meal.description}</p>
                                                    <div className="alt-meal-foods">
                                                        {meal.foods?.slice(0, 3).map((food, i) => (
                                                            <span key={i}>{food}</span>
                                                        ))}
                                                    </div>
                                                    <div className="alt-meal-meta">
                                                        <span className="alt-cal">{meal.calories}</span>
                                                        <span className="alt-cost">{meal.estimated_cost}</span>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}

                                {/* Alternative Dinner Options */}
                                {dailyMenu.dinner_alt && dailyMenu.dinner_alt.length > 0 && (
                                    <div className="alt-meals-section">
                                        <h3>🌙 Alternatif Makan Malam Lainnya</h3>
                                        <div className="alt-meals-grid">
                                            {dailyMenu.dinner_alt.map((meal, idx) => (
                                                <div key={idx} className="alt-meal-card">
                                                    <h4>{meal.title}</h4>
                                                    <p className="alt-meal-desc">{meal.description}</p>
                                                    <div className="alt-meal-foods">
                                                        {meal.foods?.slice(0, 3).map((food, i) => (
                                                            <span key={i}>{food}</span>
                                                        ))}
                                                    </div>
                                                    <div className="alt-meal-meta">
                                                        <span className="alt-cal">{meal.calories}</span>
                                                        <span className="alt-cost">{meal.estimated_cost}</span>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}

                                {dailyMenu.snacks && dailyMenu.snacks.length > 0 && (
                                    <div className="snacks-section">
                                        <h3><Apple size={20} /> Camilan Sehat</h3>
                                        <div className="snacks-grid">
                                            {dailyMenu.snacks.map((snack, idx) => (
                                                <div key={idx} className="snack-card">
                                                    <h4>{snack.title}</h4>
                                                    <div className="snack-foods">
                                                        {snack.foods?.map((food, i) => (
                                                            <span key={i}>{food}</span>
                                                        ))}
                                                    </div>
                                                    <div className="snack-meta">
                                                        <span className="snack-cal">{snack.calories}</span>
                                                        {snack.estimated_cost && (
                                                            <span className="snack-cost"><Wallet size={12} /> {snack.estimated_cost}</span>
                                                        )}
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    </div>
                                )}

                                {/* Drinks Section */}
                                {dailyMenu.drinks && dailyMenu.drinks.length > 0 && (
                                    <div className="drinks-section">
                                        <h3>🥤 Minuman Rekomendasi</h3>
                                        <div className="drinks-grid">
                                            {dailyMenu.drinks.map((drink, idx) => (
                                                <span key={idx} className="drink-item good">{drink}</span>
                                            ))}
                                        </div>
                                        {dailyMenu.avoid_drinks && dailyMenu.avoid_drinks.length > 0 && (
                                            <div className="avoid-section">
                                                <h4>Hindari:</h4>
                                                <div className="avoid-grid">
                                                    {dailyMenu.avoid_drinks.map((drink, idx) => (
                                                        <span key={idx} className="drink-item bad">{drink}</span>
                                                    ))}
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                )}

                                {/* Fruits Section */}
                                {dailyMenu.fruits && dailyMenu.fruits.length > 0 && (
                                    <div className="fruits-section">
                                        <h3>🍎 Buah Rekomendasi</h3>
                                        <div className="fruits-grid">
                                            {dailyMenu.fruits.map((fruit, idx) => (
                                                <span key={idx} className="fruit-item good">{fruit}</span>
                                            ))}
                                        </div>
                                        {dailyMenu.avoid_fruits && dailyMenu.avoid_fruits.length > 0 && (
                                            <div className="avoid-section">
                                                <h4>Hindari:</h4>
                                                <div className="avoid-grid">
                                                    {dailyMenu.avoid_fruits.map((fruit, idx) => (
                                                        <span key={idx} className="fruit-item bad">{fruit}</span>
                                                    ))}
                                                </div>
                                            </div>
                                        )}
                                    </div>
                                )}
                            </>
                        ) : (
                            <div className="empty-state">
                                <Calendar size={48} />
                                <h3>{menuError ? 'Gagal Memuat Menu' : 'Memuat Menu Harian...'}</h3>
                                <p>{menuError || 'Pastikan Anda sudah menambahkan data kesehatan untuk mendapatkan menu personal'}</p>
                                {menuError && (
                                    <button onClick={loadRecommendations} className="cta-btn" style={{ marginBottom: '10px' }}>
                                        Coba Lagi
                                    </button>
                                )}
                                <Link to="/health/add" className="cta-btn">Tambah Data Kesehatan</Link>
                            </div>
                        )}
                    </div>
                )}

                {/* Recommendations Section */}
                {(type !== 'food' || activeTab === 'recommendations') && (
                    <>
                        {recommendations.length === 0 ? (
                            <div className="empty-state">
                                <Heart size={48} />
                                <h3>Belum Ada Rekomendasi</h3>
                                <p>Lengkapi data kesehatan Anda untuk mendapatkan rekomendasi personal</p>
                                <Link to="/health/add" className="cta-btn">Tambah Data Kesehatan</Link>
                            </div>
                        ) : (
                            <div className="rec-list">
                                {recommendations.map((rec, idx) => (
                                    <div key={idx} className="rec-card">
                                        <div className="rec-header">
                                            <h3>{rec.title}</h3>
                                            <span className="rec-category">{rec.category}</span>
                                        </div>
                                        <p className="rec-desc">{rec.description}</p>
                                        <p className="rec-reason"><Lightbulb size={16} /> {rec.reason}</p>

                                        {type === 'food' && (
                                            <>
                                                {rec.foods && (
                                                    <div className="rec-section foods">
                                                        <h4><CheckCircle size={16} /> Makanan Dianjurkan</h4>
                                                        <div className="item-tags">
                                                            {rec.foods.map((food, i) => (
                                                                <span key={i} className="tag good">{food}</span>
                                                            ))}
                                                        </div>
                                                    </div>
                                                )}
                                                {rec.avoid && rec.avoid.length > 0 && (
                                                    <div className="rec-section avoid">
                                                        <h4><XCircle size={16} /> Hindari</h4>
                                                        <div className="item-tags">
                                                            {rec.avoid.map((food, i) => (
                                                                <span key={i} className="tag bad">{food}</span>
                                                            ))}
                                                        </div>
                                                    </div>
                                                )}
                                            </>
                                        )}

                                        {type === 'exercise' && (
                                            <>
                                                {rec.exercises && (
                                                    <div className="rec-section">
                                                        <h4><Flame size={16} /> Olahraga</h4>
                                                        <div className="item-tags">
                                                            {rec.exercises.map((ex, i) => (
                                                                <span key={i} className="tag good">{ex}</span>
                                                            ))}
                                                        </div>
                                                    </div>
                                                )}
                                                <div className="exercise-meta">
                                                    {rec.duration && (
                                                        <span><Clock size={14} /> {rec.duration}</span>
                                                    )}
                                                    {rec.frequency && (
                                                        <span>📅 {rec.frequency}</span>
                                                    )}
                                                    {rec.intensity && (
                                                        <span>⚡ {rec.intensity}</span>
                                                    )}
                                                </div>
                                            </>
                                        )}

                                        {type === 'emotional' && (
                                            <>
                                                {rec.activities && (
                                                    <div className="rec-section">
                                                        <h4>🎯 Aktivitas</h4>
                                                        <ul className="activity-list">
                                                            {rec.activities.map((act, i) => (
                                                                <li key={i}>{act}</li>
                                                            ))}
                                                        </ul>
                                                    </div>
                                                )}
                                                {rec.tips && (
                                                    <div className="rec-section tips">
                                                        <h4>💭 Tips</h4>
                                                        <ul className="tips-list">
                                                            {rec.tips.map((tip, i) => (
                                                                <li key={i}>{tip}</li>
                                                            ))}
                                                        </ul>
                                                    </div>
                                                )}
                                            </>
                                        )}
                                    </div>
                                ))}
                            </div>
                        )}
                    </>
                )}

                <div className="rec-nav">
                    <Link
                        to="/recommendations/food"
                        className={`nav-item ${type === 'food' ? 'active' : ''}`}
                    >
                        <Utensils size={20} />
                        <span>Makanan</span>
                    </Link>
                    <Link
                        to="/recommendations/exercise"
                        className={`nav-item ${type === 'exercise' ? 'active' : ''}`}
                    >
                        <Dumbbell size={20} />
                        <span>Olahraga</span>
                    </Link>
                    <Link
                        to="/recommendations/emotional"
                        className={`nav-item ${type === 'emotional' ? 'active' : ''}`}
                    >
                        <Heart size={20} />
                        <span>Emosional</span>
                    </Link>
                </div>
            </div>
        </div>
    );
};

export default Recommendations;
