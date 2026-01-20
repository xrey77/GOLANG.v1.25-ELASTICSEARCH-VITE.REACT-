import axios from 'axios';
import { useState, useEffect } from 'react';

const api = axios.create({
  baseURL: "http://localhost:5000",
  headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json'
  }
});

interface Product {
  id: number;
  category: string;
  descriptions: string;
  qty: number;
  unit: string;
  costprice: number;
  sellprice: number;
  productpicture: string;
  alertstocks: number;
  criticalstocks: number;
}

// Define the shape of your API response
interface ApiResponse {
  products: Product[];
  totpage: number;
  totalrecords: number;
  page: number;
}

export default function Prodlist() {
  const toDecimal = (number: number) => {
    return new Intl.NumberFormat('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(number);
  };

  const [page, setPage] = useState<number>(1);
  const [totpage, setTotpage] = useState<number>(0);
  const [totalrecs, setTotalrecs] = useState<number>(0);
  
  // Fix: Initialize as Product array instead of empty tuple []
  const [products, setProducts] = useState<Product[]>([]);

  const fetchProducts = async (pg: number) => {
    try {
      // Fix: Cast the response to your expected API structure
      const res = await api.get<ApiResponse>(`/products/list/${pg}`);
      
      // Update states based on response structure
      setProducts(res.data.products); 
      setTotpage(res.data.totpage);
      setTotalrecs(res.data.totalrecords);
      setPage(res.data.page);
    } catch (error: any) {
      console.error(error.response?.data?.message || "An error occurred");
    }
  };

  useEffect(() => {
    fetchProducts(page);
  }, [page]);

  const firstPage = (event: React.MouseEvent) => {
    event.preventDefault();
    setPage(1);
  };

  const nextPage = (event: React.MouseEvent) => {
    event.preventDefault();
    if (page < totpage) {
      setPage(prev => prev + 1);
    }
  };

  const prevPage = (event: React.MouseEvent) => {
    event.preventDefault();
    if (page > 1) {
      setPage(prev => prev - 1);
    }
  };

  const lastPage = (event: React.MouseEvent) => {
    event.preventDefault();
    setPage(totpage);
  };

  return (
    <div className="container">
      <h1 className='text-warning embossed mt-3'>Products List</h1>
      <table className="table table-danger table-striped">
        <thead>
          <tr>
            <th className="bg-primary text-white">#</th>
            <th className="bg-primary text-white">Descriptions</th>
            <th className="bg-primary text-white">Qty</th>
            <th className="bg-primary text-white">Unit</th>
            <th className="bg-primary text-white">Price</th>
          </tr>
        </thead>
        <tbody>
          {products.map((item) => (
            <tr key={item.id}>
              <td>{item.id}</td>
              <td>{item.descriptions}</td>
              <td>{item.qty}</td>
              <td>{item.unit}</td>
              <td>₱{toDecimal(item.sellprice)}</td>
            </tr>
          ))}
        </tbody>
      </table>

      <nav aria-label="Page navigation example">
        <ul className="pagination sm">
          <li className="page-item"><a onClick={firstPage} className="page-link sm" href="#">First</a></li>
          <li className="page-item"><a onClick={prevPage} className="page-link sm" href="#">Previous</a></li>
          <li className="page-item"><a onClick={nextPage} className="page-link sm" href="#">Next</a></li>
          <li className="page-item"><a onClick={lastPage} className="page-link sm" href="#">Last</a></li>
          <li className="page-item page-link text-danger sm">Page {page} of {totpage}</li>
        </ul>
      </nav>
      <div className='text-warning'><strong>Total Records : {totalrecs}</strong></div>
    </div>
  );
}
