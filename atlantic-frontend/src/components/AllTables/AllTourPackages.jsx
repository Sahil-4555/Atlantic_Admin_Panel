import * as React from "react";
import { Link, Navigate } from "react-router-dom";
import axios from "axios";
import { useState, useEffect } from "react";
import "./AllTables.scss";

const AllTourPackages = () => {
  const [getTourPackageDetails, setTourPackageDetails] = useState([]);

  useEffect(() => {
    getTourPackage();
  }, []);

  const getTourPackage = () => {
    axios
      .get("http://localhost:5000/allproducts")
      .then((response) => response.data)
      .then((data) => {
        setTourPackageDetails(data.product);
      });
  };

  const handleView = (_id) => {
    window.localStorage.setItem("tourSpecificId", _id);
  };

  return window.localStorage.getItem("token") ? (
    <div className="allproducts">
      {getTourPackageDetails.map((details) => (
        <div className="card">
          <img
            // src={details.image}
            className="card-img-top"
          />
          <div className="card-body" style={{paddingLeft:"50px"}}>
            <h5 className="card-title">{details.title}</h5>
            <p className="card-text">{details.description}</p>
            <h6>Color : {details.color}</h6>
            <h6>Size : {details.size}</h6>
            <div className="rowdatas">
              <div className="product-card-price">
                <span>
                  <del>{details.price} ₹</del>
                </span>
                {/* <span className="curr-price">{details.priceDiscount} ₹</span> */}
              </div>
              <Link
                // to={"/tourpackages/" + details._id}
                onClick={() => handleView(details.productid)}
                style={{ textDecoration: "none" }}
                className="btn DetailsButton"
              >
                View Details
              </Link>
            </div>
          </div>
        </div>
      ))}
    </div>
  ) : (
    <Navigate to="/login" />
  );
};

export default AllTourPackages;
