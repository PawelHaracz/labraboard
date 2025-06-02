import React from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import {Box, Button, useTheme} from "@mui/material";
import {tokens} from "../../theme";
import Typography from "@mui/material/Typography";

const ProjectDetails = () => {
    const { id } = useParams();
    const theme = useTheme();
    const colors = tokens(theme.palette.mode);
    const navigate = useNavigate();

    const handleGoBack = () => {
        navigate(-1); // This navigates back to the previous page
    };

    return (
        <Box m="20px">
            <Button onClick={handleGoBack} variant="contained" color="primary">
                Go Back
            </Button>
            <Typography variant="h4" component="h1" mt={2}>
                Project Details
            </Typography>
            <Typography variant="body1" component="p" mt={1}>
                Project ID: {id}
            </Typography>
            {/* Fetch and display more details about the project using the id */}
        </Box>
    );
}

export default ProjectDetails;