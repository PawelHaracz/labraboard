import {Box, Button, Grid, useTheme} from "@mui/material";
import {tokens} from "../../theme";
import Header from "../../components/Header";
import AddCircleOutlineOutlinedIcon from '@mui/icons-material/AddCircleOutlineOutlined';
import {mockDataProjects} from "../../data/mockProjects";
import ProjectCard from "../../components/ProjectCard"

const Projects = () => {
    const theme = useTheme();
    const colors = tokens(theme.palette.mode);
    return (
        <Box m="20px">
            {/* HEADER */}
            <Box display="flex" justifyContent="space-between" alignItems="center">
                <Header title="Projects" subtitle="Your projects" />

                <Box>
                    <Button
                        sx={{
                            backgroundColor: colors.greenAccent[700],
                            color: colors.grey[100],
                            fontSize: "14px",
                            fontWeight: "bold",
                            padding: "10px 20px",
                        }}
                    >
                        <AddCircleOutlineOutlinedIcon sx={{ mr: "10px" }} />
                        Add New Project
                    </Button>
                </Box>

            </Box>
                <Grid container spacing={2} mt={2}>
                    {mockDataProjects.map((project) => (
                        <Grid item xs={12} sm={6} md={4} key={project.id}>
                            <ProjectCard data={project} />
                        </Grid>
                    ))}
                </Grid>

        </Box>
    )
}
export default Projects;