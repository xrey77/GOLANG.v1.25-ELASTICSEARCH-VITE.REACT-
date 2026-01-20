import React, { useEffect, useState } from 'react';
import axios from 'axios';

const api = axios.create({
  baseURL: "http://localhost:5000",
  headers: {
    'Accept': 'application/pdf',
    'Content-Type': 'application/pdf'
  }
});

const ProductReport: React.FC = () => {
  const [loading, setLoading] = useState<boolean>(false);
  const [pdfUrl, setPdfUrl] = useState<string | null>(null);

  const fetchPdf = async () => {
    setLoading(true);
    try {
      const response = await api.get('/productreport', {
        responseType: 'blob', 
      });

      // Create a blob from the response
      const blob = new Blob([response.data], { type: 'application/pdf' });
      
      // Create a URL for the blob
      const url = window.URL.createObjectURL(blob);
      setPdfUrl(url);
    } catch (error) {
      console.error('Failed to load PDF:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPdf();
    
    return () => {
      if (pdfUrl) window.URL.revokeObjectURL(pdfUrl);
    };
  }, []);

  return (
    <div className="container-fluid">      
      {loading && <p>Loading PDF...</p>}

      {pdfUrl && (
        <div className="container-fluid">
          <iframe
            src={`${pdfUrl}#toolbar=1`} // #toolbar=0 is optional to hide browser controls
            width="100%"
            height="800px"
            title="Product Report"
          />
        </div>
      )}
    </div>
  );
};

export default ProductReport;
