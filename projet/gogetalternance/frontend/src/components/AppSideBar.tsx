import { Home, Search, Star, User } from "lucide-react";
import { NavLink } from "@/components/NavLink";
import { useLocation } from "react-router-dom";
import {
    Sidebar,
    SidebarContent,
    SidebarGroup,
    SidebarGroupContent,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarFooter,
    useSidebar,
} from "./ui/sidebar.tsx";

const navItems = [
    { title: "Home", url: "/", icon: Home },
    { title: "Search", url: "/search", icon: Search },
    { title: "Favorites", url: "/favorites", icon: Star },
];

interface AppSidebarProps {
    onAuthClick: () => void;
}

export function AppSidebar({ onAuthClick }: AppSidebarProps) {
    const { state } = useSidebar();
    const collapsed = state === "collapsed";
    const location = useLocation();

    return (
        <Sidebar collapsible="icon">
            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroupContent>
                        <SidebarMenu>
                            {!collapsed && (
                                <div className="px-3 py-4 mb-2">
                                    <h2 className="text-lg font-bold tracking-tight">GoGet<span className="text-muted-foreground">Alternance</span></h2>
                                </div>
                            )}
                            {navItems.map((item) => (
                                <SidebarMenuItem key={item.title}>
                                    <SidebarMenuButton asChild>
                                        <NavLink
                                            to={item.url}
                                            end={item.url === "/"}
                                            className="hover:bg-accent"
                                            activeClassName="bg-accent font-medium"
                                        >
                                            <item.icon className="mr-2 h-4 w-4" />
                                            {!collapsed && <span>{item.title}</span>}
                                        </NavLink>
                                    </SidebarMenuButton>
                                </SidebarMenuItem>
                            ))}
                        </SidebarMenu>
                    </SidebarGroupContent>
                </SidebarGroup>
            </SidebarContent>
            <SidebarFooter>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <SidebarMenuButton onClick={onAuthClick} className="hover:bg-accent">
                            <User className="mr-2 h-4 w-4" />
                            {!collapsed && <span>Account</span>}
                        </SidebarMenuButton>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarFooter>
        </Sidebar>
    );
}
