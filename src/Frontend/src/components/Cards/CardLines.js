import React, { useState } from "react";
import { Table, TableHead, TableRow, TableCell, TableBody } from '@mui/material';
import Pagination from '@mui/material/Pagination';


const itemsPerPage = 5;

export default function CardLines({ onClick, data
}) {
    const [currentPage, setCurrentPage] = useState(1);

    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;

    const currentData = data.slice(startIndex, endIndex);

    const totalPages = Math.ceil(data.length / itemsPerPage);

    const handlePageChange = (newPage) => {
        setCurrentPage(newPage);
    };

    return (
        <>
            <div className="relative flex flex-col min-w-0 break-words bg-white rounded mb-6 xl:mb-0 shadow-lg">
                <div className="flex-auto p-4">
                    <div className="flex flex-wrap">
                        <div className="relative w-full pr-4 max-w-full flex-grow flex-1">
                            <span className="font-semibold text-xl text-blueGray-700">
                                Memoria
                            </span>
                        </div>
                        <button
                            className="bg-indigo-500 text-white active:bg-indigo-600 text-xs font-bold uppercase px-3 py-1 rounded outline-none focus:outline-none mr-1 mb-1 ease-linear transition-all duration-150"
                            type="button"
                            style={{ backgroundColor: "#271744" }}
                            onClick={onClick}
                        >
                            <i className="fas fa-play text-white mr-4"></i>
                            Start
                        </button>

                    </div>
                    <div className="flex flex-wrap" style={{ paddingTop: 10 }}>
                        <div style={{ height: 250, width: 545 }}>

                            <div>


                                <div>
                                    <Table size="small">
                                        <TableHead>
                                            <TableRow>
                                                <TableCell style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000', fontWeight: 'bold' }}>Address</TableCell>
                                                <TableCell style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000', fontWeight: 'bold' }}>Data</TableCell>
                                            </TableRow>
                                        </TableHead>
                                        <TableBody>
                                            {currentData.map((item, index) => (
                                                <TableRow key={index}>
                                                    <TableCell style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000' }}>{item.Address}</TableCell>
                                                    <TableCell style={{ borderColor: '#000', border: 'solid', borderWidth: 0.5, backgroundColor: '#e2e8f0', color: '#000' }}>{item.Data}</TableCell>
                                                </TableRow>
                                            ))}
                                        </TableBody>
                                    </Table>

                                    <div className="pagination" style={{ display: 'flex', justifyContent: 'center' }}>
                                        <Pagination
                                            count={totalPages}
                                            page={currentPage}
                                            onChange={(event, page) => handlePageChange(page)}
                                            variant="outlined"
                                            shape="rounded"
                                            size="small"
                                        />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </>
    );
}
