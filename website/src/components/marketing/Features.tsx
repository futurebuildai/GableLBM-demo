import { motion } from "framer-motion";
import { Search, MapPin, Camera, MousePointer2 } from "lucide-react";

export const Features = () => {
    const features = [
        {
            title: "Field Visibility",
            subtitle: "Logistics",
            description: "Live truck tracking for your customers. Strip 40% of 'where is it?' calls out of the dispatch office while building contractor trust.",
            icon: <MapPin className="w-5 h-5" />,
            image: "https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?auto=format&fit=crop&q=80&w=800",
        },
        {
            title: "Instant SKU Discovery",
            subtitle: "Counter Ops",
            description: "Fuzzy search that understands lumber yard terminology. New counter hires find the right stud and SKU in minutes, not months.",
            icon: <Search className="w-5 h-5" />,
            image: "https://images.unsplash.com/photo-1589939705384-5185137a7f0f?auto=format&fit=crop&q=80&w=800",
        },
        {
            title: "Visual Proof of Delivery",
            subtitle: "Liability Control",
            description: "High-res drop photos attached directly to the ticket. Kill 90% of damage and shortage claims before they start.",
            icon: <Camera className="w-5 h-5" />,
            image: "https://images.unsplash.com/photo-1517581177682-a085bb7ffb15?auto=format&fit=crop&q=80&w=800",
        },
        {
            title: "Yard-Ready Interface",
            subtitle: "User Experience",
            description: "Low-friction UI built for the high-pressure environment of the yard. Zero training required for the digital native.",
            icon: <MousePointer2 className="w-5 h-5" />,
            image: "https://images.unsplash.com/photo-1550751827-4bd374c3f58b?auto=format&fit=crop&q=80&w=800",
        },
    ];

    return (
        <section className="py-24 px-4">
            <div className="container mx-auto">
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center mb-24">
                    <div>
                        <h2 className="text-4xl md:text-5xl font-bold mb-6">
                            Consumer Grade <br />
                            <span className="text-gable-blue">Enterprise Power.</span>
                        </h2>
                        <p className="text-xl text-slate-400 font-light max-w-xl">
                            We've mapped modern logic into the pro-dealer environment. Gable equips independent yards with the tech scale required to compete with the giants.
                        </p>
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                        {/* Dynamic layout elements could go here */}
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                    {features.map((feature, i) => (
                        <motion.div
                            key={i}
                            initial={{ opacity: 0, scale: 0.98 }}
                            whileInView={{ opacity: 1, scale: 1 }}
                            viewport={{ once: true }}
                            transition={{ delay: i * 0.1 }}
                            className="glass-card overflow-hidden group hover-depth"
                        >
                            <div className="flex flex-col md:flex-row h-full">
                                <div className="flex-1 p-8">
                                    <div className="flex items-center space-x-3 text-gable-blue mb-4">
                                        {feature.icon}
                                        <span className="text-sm font-mono uppercase tracking-widest">{feature.subtitle}</span>
                                    </div>
                                    <h3 className="text-2xl font-bold mb-4">{feature.title}</h3>
                                    <p className="text-slate-400 font-light leading-relaxed">
                                        {feature.description}
                                    </p>
                                </div>
                                <div className="md:w-1/3 relative overflow-hidden h-48 md:h-auto">
                                    <img
                                        src={feature.image}
                                        alt={feature.title}
                                        className="absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                                    />
                                    <div className="absolute inset-0 bg-gable-bg/40 group-hover:bg-transparent transition-colors duration-500" />
                                </div>
                            </div>
                        </motion.div>
                    ))}
                </div>
            </div>
        </section>
    );
};
