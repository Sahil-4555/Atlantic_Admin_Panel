import Home from "./pages/home/Home";
import Login from "./pages/login/Login";
import UserList from "./pages/list/UserList";
import SpecificUserDetail from "./pages/SpecificItemDetails/SpecificUserDetail";
import SpecificUserDetailEdit from "./pages/EditItems/SpecificUserDetailEdit";
import AddingUser from "./pages/AdingsItems/AddingUser";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { productInputs } from "./formSource";
import "./style/dark.scss";
import { useContext } from "react";
import { DarkModeContext } from "./context/darkModeContext";
import Transaction from "./pages/Transactions/Transaction";
import TourPacakageList from "./pages/list/TourPackageList";

function App() {
  const { darkMode } = useContext(DarkModeContext);

  return (
    <div className={darkMode ? "app dark" : "app"}>
      <BrowserRouter>
        <Routes>
          <Route path="/">
            {/* Home Component */}
            <Route index element={<Home />} />
            {/* Login Component */}
            <Route path="login" element={<Login />} />

            {/* Users Related Component */}

            <Route path="allusers">
              {/* Main Indec Componen */}
              <Route index element={<UserList />} />

              {/* SpecificUserDetail Component */}
              <Route path=":userId" element={<SpecificUserDetail />} />

              {/* SpecificUserEdit Component */}
              <Route path="edit/:userId" element={<SpecificUserDetailEdit />} />

              {/* Adding New User */}
              {/* <Route
                path="newuser"
                // element={<New inputs={userInputs} title="Add New User" />}
                element={<AddingUser title="Add New User" />}
              /> */}


            </Route>


            {/* TourPackages Related Component */}
            <Route path="allproduct">
              {/* Main Index Component */}
              <Route index element={<TourPacakageList />} />
              {/* SpecificTourPackagesDetail Component */}
              {/* <Route path=":tourpackagesId" element={<SpecificUserDetail />} /> */}
            </Route>
            {/* Adding New TourPackage */}
            {/* <Route
              path="newtourpackages"
              element={<AddingUser title="Add New TourPackages" />}
            /> */}

            {/* <Route path="transactions" element={<Transaction />} /> */}
          </Route>
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;