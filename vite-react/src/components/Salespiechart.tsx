import axios from 'axios';
import { useState, useEffect } from 'react';

const api = axios.create({
  baseURL: "http://localhost:5000",
  headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json'
  }
});

export default function Salespiechart() {
const [chartUrl, setChartUrl] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchChart = async () => {
    try {
      const response = await api.get('/sales/piechart', {
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
    <div className="container-fluid bg-white mt-3">
      {chartUrl ? (
        <div className="d-flex justify-content-center">
        <img className="img-fluid"            
          src={chartUrl} 
          alt="Monthly Sales Chart"           
        />
        </div>
      ) : (
        <p>No chart data available.</p>
      )}
    </div>    
    );
}
