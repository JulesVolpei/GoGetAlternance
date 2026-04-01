import { NavLink as RouterNavLink } from "react-router-dom";
import { forwardRef, ComponentProps } from "react"; // Ajout de ComponentProps
import { cn } from "@/lib/utils";

// On utilise ComponentProps<typeof RouterNavLink> à la place de NavLinkProps
interface NavLinkCompatProps extends Omit<ComponentProps<typeof RouterNavLink>, "className"> {
    className?: string;
    activeClassName?: string;
    pendingClassName?: string;
}

const NavLink = forwardRef<HTMLAnchorElement, NavLinkCompatProps>(
    ({ className, activeClassName, pendingClassName, to, ...props }, ref) => {
        return (
            <RouterNavLink
                ref={ref}
                to={to}
                className={({ isActive, isPending }) =>
                    cn(className, isActive && activeClassName, isPending && pendingClassName)
                }
                {...props}
            />
        );
    },
);

NavLink.displayName = "NavLink";

export { NavLink };