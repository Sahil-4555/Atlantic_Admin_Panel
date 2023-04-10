import "./AllTables.scss";
import { DataGrid } from "@mui/x-data-grid";
import { Link, Navigate } from "react-router-dom";
import axios from "axios";
import { useState, useEffect } from "react";

const AllUsers = () => {
  // const [data, setData] = useState(userRows);
  const [getUserDetails, setUserDetails] = useState([]);

  useEffect(() => {
    getUsers();
  }, []);

  const getUsers = () => {
    axios
      .get("http://localhost:5000/allusers")
      .then((response) => response.data)
      .then((data) => {
        setUserDetails(data.user);
        // for (var i = 0; i < data.data.Users.length; i++) {
        //   console.log(data.data.Users);
        // }
      });
  };

  const handleDelete = (_id) => {
    axios.delete("http://localhost:5000/users/" + _id).then((response) => {
      // if (response.data !== null) {
      //   setUserDetails(getUserDetails.filter((item) => item._id !== _id));
      //   //console.log(response.data.Users.username + "Record Delete Successful");
      // }
      window.location.href = 'http://localhost:3000/allusers'
    });
  };

  const handleView = (_id) => {
    window.localStorage.setItem("userSpecificId", _id);
  };

  const actionColumn = [
    {
      field: "action",
      headerName: "Action",
      width: 200,
      renderCell: (params) => {
        return (
          <div className="cellAction">
            <Link
              to={"/allusers/" + params.row.uid}
              onClick={() => handleView(params.row.uid)}
              style={{ textDecoration: "none" }}
            >
              <div className="viewButton">View</div>
            </Link>
            <div
              className="deleteButton"
              onClick={() => handleDelete(params.row.uid)}
            >
              Delete
            </div>
          </div>
        );
      },
    },
  ];
  return window.localStorage.getItem("token") ? (
    <div className="datatable">
      {/* <div className="datatableTitle">
        Add New User
        <Link to="newuser" className="link">
          Add New
        </Link>
      </div> */}
      <DataGrid
        className="datagrid"
        rows={getUserDetails}
        columns={userColumns.concat(actionColumn)}
        pageSize={8}
        getRowId={(row) => row.id}
        rowsPerPageOptions={[8]}
        checkboxSelection
        onClick={getUsers}
      />
    </div>
  ) : (
    <Navigate to="/login" />
  );
};

export const userColumns = [
  {
    field: "id",
    headerName: "ID",
    width: 150,
  },
  {
    field: "photourl",
    headerName: "Profile",
    width: 100,
    renderCell: (params) => {
      return (
        <div className="cellWithImg">
          <img className="cellImg" src={params.row.photourl} alt="" />
        </div>
      );
    },
  },
  {
    field: "name",
    headerName: "Users",
    width: 200,
  },
  {
    field: "email",
    headerName: "Email",
    width: 280,
  },
];

export default AllUsers;
