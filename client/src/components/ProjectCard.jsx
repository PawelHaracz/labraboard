import React from 'react';
import {Card, CardContent, CardActions, Typography, Button, Box, Chip, useTheme} from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import {tokens} from "../theme";
import {Link} from "react-router-dom";

const ProjectCard = ({ data }) => {
    const theme = useTheme();
    const colors = tokens(theme.palette.mode);
    return (
        <Card sx={{
            minWidth: 275,
            backgroundColor: colors.primary[900],
            color: colors.grey[100],
            fontSize: "15px",
            fontWeight: "bold",
            padding: "10px 20px",
        }}>
            <CardContent>
                <Link to={`/project/${data.id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
                    <Typography variant="h5" component="div">
                        {data.name}
                    </Typography>
                </Link>
                {/*<Box display="flex" alignItems="center" my={1}>*/}
                {/*    <CheckCircleIcon color="success" />*/}
                {/*    <Typography variant="body2" color="success.main" sx={{ ml: 1 }}>*/}
                {/*        Healthy | Synced*/}
                {/*    </Typography>*/}
                {/*</Box>*/}
                <Typography>
                    Repository: <a href={data.repositoryUrl} target="_blank" rel="noopener noreferrer">{data.repositoryUrl}</a>
                </Typography>
                <Typography>
                    Target Revision: {data.repositoryBranch}
                </Typography>
                <Typography>
                    Path: {data.repositoryTerraformPath}
                </Typography>
            </CardContent>
            <CardActions>
                <Button sx={{
                    backgroundColor: colors.primary[900],
                    color: colors.grey[100],
                    fontSize: "15px",
                    fontWeight: "bold",
                    padding: "10px 20px",
                }}>SYNC</Button>
                <Button sx={{
                    backgroundColor: colors.primary[900],
                    color: colors.grey[100],
                    fontSize: "15px",
                    fontWeight: "bold",
                    padding: "10px 20px",
                }}>DELETE</Button>
            </CardActions>
        </Card>
    );
}
export default ProjectCard