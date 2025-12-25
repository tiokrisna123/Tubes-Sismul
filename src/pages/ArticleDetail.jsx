import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { articlesAPI } from '../services/api';
import { ArrowLeft, Clock, Loader2 } from 'lucide-react';
import './ArticleDetail.css';

const ArticleDetail = () => {
    const { id } = useParams();
    const navigate = useNavigate();
    const [article, setArticle] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetchArticle();
    }, [id]);

    const fetchArticle = async () => {
        try {
            const res = await articlesAPI.getById(id);
            setArticle(res.data);
        } catch (error) {
            console.error('Failed to fetch article:', error);
        } finally {
            setLoading(false);
        }
    };

    const getCategoryIcon = (cat) => {
        const icons = {
            'nutrisi': '🥗',
            'olahraga': '🏃',
            'mental': '🧠',
            'tidur': '😴',
            'umum': '❤️'
        };
        return icons[cat] || '📄';
    };

    if (loading) {
        return (
            <div className="article-detail-loading">
                <Loader2 className="spinner" size={48} />
                <p>Memuat artikel...</p>
            </div>
        );
    }

    if (!article) {
        return (
            <div className="article-detail-loading">
                <p>Artikel tidak ditemukan</p>
                <button onClick={() => navigate('/articles')}>Kembali</button>
            </div>
        );
    }

    return (
        <div className="article-detail-page">
            <header className="article-detail-header">
                <button className="back-btn" onClick={() => navigate('/articles')}>
                    <ArrowLeft size={24} />
                </button>
                <span className="article-category-badge">
                    {getCategoryIcon(article.category)} {article.category}
                </span>
            </header>

            {article.image_url && (
                <div className="article-hero">
                    <img src={article.image_url} alt={article.title} />
                </div>
            )}

            <article className="article-detail-content">
                <h1>{article.title}</h1>

                <div className="article-detail-meta">
                    <span className="read-time">
                        <Clock size={16} />
                        {article.read_time} menit baca
                    </span>
                </div>

                <div className="article-body">
                    {article.content.split('\n\n').map((paragraph, index) => {
                        if (paragraph.startsWith('## ')) {
                            return <h2 key={index}>{paragraph.replace('## ', '')}</h2>;
                        } else if (paragraph.startsWith('- ')) {
                            return (
                                <ul key={index}>
                                    {paragraph.split('\n').map((item, i) => (
                                        <li key={i}>{item.replace('- ', '')}</li>
                                    ))}
                                </ul>
                            );
                        } else if (paragraph.includes('**')) {
                            return (
                                <p key={index} dangerouslySetInnerHTML={{
                                    __html: paragraph.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
                                }} />
                            );
                        } else {
                            return <p key={index}>{paragraph}</p>;
                        }
                    })}
                </div>
            </article>
        </div>
    );
};

export default ArticleDetail;
