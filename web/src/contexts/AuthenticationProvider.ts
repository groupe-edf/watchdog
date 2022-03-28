import React from "react"

const AuthContext = React.createContext(null)

export const AuthProvider = ({ userData, children }: any) => {
}

export const useAuth = () => React.useContext(AuthContext)
