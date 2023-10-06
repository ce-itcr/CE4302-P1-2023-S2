import React from "react";
import { Link } from "react-router-dom";


export default function Header({ title, onResume, onClose, onDownload }) {

    const checkDashboard = () => {
        if (title === '/app/dashboard') {
            return (<li className="flex items-center">
                <button
                    className="bg-lightBlue-500 text-white active:bg-lightBlue-600 text-xs font-bold uppercase px-4 py-2 rounded shadow hover:shadow-lg outline-none focus:outline-none lg:mr-1 lg:mb-0 ml-3 mb-3 ease-linear transition-all duration-150"
                    type="button"
                    style={{ backgroundColor: "#271744" }}
                    onClick={onResume}
                >
                    <i className="fas fa-chart-simple"></i> Resumen
                </button>
            </li>)
        }
    }

    const checkSummary = () => {
        if (title === '/app/summary') {
            return (<>
            {/* <li className="flex items-center">
                <button
                    className="bg-lightBlue-500 text-white active:bg-lightBlue-600 text-xs font-bold uppercase px-4 py-2 rounded shadow hover:shadow-lg outline-none focus:outline-none lg:mr-1 lg:mb-0 ml-3 mb-3 ease-linear transition-all duration-150"
                    type="button"
                    style={{ backgroundColor: "#271744" }}
                    onClick={onClose}
                >
                    <i className="fas fa-arrow-alt-circle-down"></i> Descargar
                </button>
            </li> */}
            <li className="flex items-center">
                    <button
                        className="bg-lightBlue-500 text-white active:bg-lightBlue-600 text-xs font-bold uppercase px-4 py-2 rounded shadow hover:shadow-lg outline-none focus:outline-none lg:mr-1 lg:mb-0 ml-3 mb-3 ease-linear transition-all duration-150"
                        type="button"
                        style={{ backgroundColor: "#271744" }}
                        onClick={onClose}
                    >
                        <i className="fas fa-close"></i> Cerrar
                    </button>
                </li></>)
        }
    }

    return (
        <>
            <nav className="top-0 fixed z-50 w-full flex flex-wrap items-center justify-between px-2 py-3 navbar-expand-lg bg-white shadow">
                <div className="container px-4 mx-auto flex flex-wrap items-center justify-between">
                    <div className="w-full relative flex justify-between lg:w-auto lg:static lg:block lg:justify-start">
                        <Link
                            className="text-blueGray-700 text-sm font-bold leading-relaxed inline-block mr-4 py-2 whitespace-nowrap uppercase"
                        >
                            <img
                                className=""
                                style={{ width: '70%' }}
                                src={require("../../assets/img/Logo.png")}
                                alt="..."
                            />
                        </Link>
                    </div>
                    <ul className="flex flex-col lg:flex-row list-none lg:ml-auto">
                        <li className="flex items-center">
                            <div style={{ fontWeight: "bold", color: '#271744' }}>{title}</div>

                        </li>
                        {checkDashboard()}
                        {checkSummary()}
                    </ul>
                </div>

            </nav>
        </>
    );
}
