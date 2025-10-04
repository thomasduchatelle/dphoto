import {useEffect} from "react";
import {useNavigate} from "react-router-dom";

export default function IndexPage() {
    const navigate = useNavigate();

    useEffect(() => {
        navigate('/albums');
    }, [navigate]);

    return null;
}
