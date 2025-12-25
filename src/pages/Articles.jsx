import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { articlesAPI } from '../services/api';
import { ArrowLeft, Search, BookOpen, Clock, Loader2 } from 'lucide-react';
import './Articles.css';

const Articles = () => {
    const navigate = useNavigate();
    const [articles, setArticles] = useState([]);
    const [categories, setCategories] = useState([]);
    const [selectedCategory, setSelectedCategory] = useState('all');
    const [searchQuery, setSearchQuery] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchCategories();
        fetchArticles('all');
    }, []);

    const fetchCategories = async () => {
        try {
            const res = await articlesAPI.getCategories();
            setCategories(res.data);
        } catch (error) {
            console.error('Failed to fetch categories:', error);
        }
    };

    const fetchArticles = async (category) => {
        setLoading(true);
        try {
            const res = await articlesAPI.getAll(category === 'all' ? '' : category);
            setArticles(res.data || []);
        } catch (error) {
            console.error('Failed to fetch articles:', error);
        } finally {
            setLoading(false);
        }
    };

    const handleCategoryChange = (category) => {
        setSelectedCategory(category);
        setSearchQuery('');
        fetchArticles(category);
    };

    const handleSearch = async (e) => {
        e.preventDefault();
        if (!searchQuery.trim()) {
            fetchArticles(selectedCategory);
            return;
        }
        setLoading(true);
        try {
            const res = await articlesAPI.search(searchQuery);
            setArticles(res.data || []);
        } catch (error) {
            console.error('Search failed:', error);
        } finally {
            setLoading(false);
        }
    };

    const getCategoryIcon = (cat) => {
        const icons = {
            'all': '📚',
            'nutrisi': '🥗',
            'olahraga': '🏃',
            'mental': '🧠',
            'tidur': '😴',
            'umum': '❤️'
        };
        return icons[cat] || '📄';
    };

    if (loading && articles.length === 0) {
        return (
            <div className="articles-loading">
                <Loader2 className="spinner" size={48} />
                <p>Memuat artikel...</p>
            </div>
        );
    }

    return (
        <div className="articles-page">
            <header className="articles-header">
                <button className="back-btn" onClick={() => navigate('/dashboard')}>
                    <ArrowLeft size={24} />
                </button>
                <h1>📚 Tips & Artikel Kesehatan</h1>
            </header>

            <div className="articles-container">
                {/* Search Bar */}
                <form className="search-bar" onSubmit={handleSearch}>
                    <Search size={20} />
                    <input
                        type="text"
                        placeholder="Cari artikel..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                    />
                    <button type="submit">Cari</button>
                </form>

                {/* Categories */}
                <div className="categories-scroll">
                    {categories.map((cat) => (
                        <button
                            key={cat.id}
                            className={`category-chip ${selectedCategory === cat.id ? 'active' : ''}`}
                            onClick={() => handleCategoryChange(cat.id)}
                        >
                            <span>{cat.icon}</span>
                            <span>{cat.name}</span>
                        </button>
                    ))}
                </div>

                {/* Articles List */}
                {loading ? (
                    <div className="articles-loading-inline">
                        <Loader2 className="spinner" size={32} />
                    </div>
                ) : articles.length === 0 ? (
                    <div className="empty-state">
                        <BookOpen size={64} />
                        <h3>Belum Ada Artikel</h3>
                        <p>Tidak ada artikel yang ditemukan untuk kategori ini.</p>
                    </div>
                ) : (
                    <div className="articles-grid">
                        {articles.map((article) => (
                            <article
                                key={article.id}
                                className="article-card"
                                onClick={() => navigate(`/articles/${article.id}`)}
                            >
                                {article.image_url && (
                                    <div className="article-image">
                                        <img src={article.image_url} alt={article.title} />
                                    </div>
                                )}
                                <div className="article-content">
                                    <span className="article-category">
                                        {getCategoryIcon(article.category)} {article.category}
                                    </span>
                                    <h3>{article.title}</h3>
                                    <p>{article.summary}</p>
                                    <div className="article-meta">
                                        <Clock size={14} />
                                        <span>{article.read_time} menit baca</span>
                                    </div>
                                </div>
                            </article>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};

export default Articles;
