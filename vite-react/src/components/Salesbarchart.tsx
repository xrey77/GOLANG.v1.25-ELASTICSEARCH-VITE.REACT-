import axios from 'axios';
import { useState, useEffect } from 'react';

const api = axios.create({
  baseURL: "http://localhost:5000",
  headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json'
  }
});


export default function Salesbarchart() {
const [chartUrl, setChartUrl] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchChart = async () => {
    try {
      const response = await api.get('/sales/barchart', {
        responseType: 'blob', 
      });

      // Create a blob from the response
      const blob = new Blob([response.data], { type: 'application/pdf' });
      const objectUrl = URL.createObjectURL(blob);
      setChartUrl(objectUrl);      
      
    } catch (error) {
      console.error('Failed to load PDF:', error);
    } finally {
      setLoading(false);
    }
  };

    fetchChart();

    return () => {
      if (chartUrl) URL.revokeObjectURL(chartUrl);
    };
  }, []);

  if (loading) return <div>Loading sales data...</div>;

  return (
    <div style={{ textAlign: 'center', padding: '20px' }}>
      {chartUrl ? (
        <img 
          src={chartUrl} 
          alt="Monthly Sales Chart" 
          style={{ maxWidth: '100%', height: 'auto', border: '1px solid #ddd' }} 
        />
      ) : (
        <p>No chart data available.</p>
      )}
    </div>
    );
}
