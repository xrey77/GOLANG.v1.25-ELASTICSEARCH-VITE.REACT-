import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import { HashRouter, Route, Routes } from 'react-router-dom'
import Header from './components/Header.tsx'
import About from './components/About.tsx'
import Home from './components/Home.tsx'
import Contact from './components/Contact.tsx';
import Prodlist from './components/Prodlist.tsx';
import Prodcatalog from './components/Prodcatalog.tsx';
import Prodsearch from './components/Prodsearch.tsx';
import Profile from './components/Profile.tsx';
import Productreport from './components/Productreport.tsx';
import Salesbarchart from './components/Salesbarchart.tsx';
import Salespiechart from './components/Salespiechart.tsx';
import './App.css'

function App() {

  return (
    <HashRouter>
      <Header/>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/about" element={<About />} />
          <Route path="/contact" element={<Contact />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/productlist" element={<Prodlist />} />
          <Route path="/productcatalog" element={<Prodcatalog />} />
          <Route path="/productsearch" element={<Prodsearch />} />
          <Route path="/pdfreport" element={<Productreport />} />
          <Route path="/salesbarchart" element={<Salesbarchart />} />
          <Route path="/salespiechart" element={<Salespiechart />} />

        </Routes>
    </HashRouter>    
  )
}

export default App
