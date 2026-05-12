mod core;
mod adapters;
mod ui;

use std::sync::Arc;
use eframe::egui;
use tokio::runtime::Runtime;

use crate::core::application::SessionService;
use crate::core::domain::Session;
use crate::adapters::storage::MemoryStorageAdapter;

fn main() -> eframe::Result<()> {
    // Create tokio runtime for async
    let rt = Runtime::new().unwrap();

    // Create storage adapter (hexagonal: adapter implements port)
    let storage = Arc::new(MemoryStorageAdapter::new());

    // Create service (hexagonal: depends on port, not adapter)
    let service = Arc::new(SessionService::new(storage));

    let options = eframe::NativeOptions {
        viewport: egui::ViewportBuilder::default()
            .with_inner_size([500.0, 450.0])
            .with_title("CannaNote"),
        ..Default::default()
    };

    eframe::run_native(
        "CannaNote",
        options,
        Box::new(move |_cc| Box::new(CannaNote::new(rt, service))),
    )
}

struct CannaNote {
    rt: Runtime,
    service: Arc<SessionService>,
    // Form state
    method: String,
    amount: String,
    unit: String,
    // Session list
    sessions: Vec<Session>,
    status_message: String,
}

impl CannaNote {
    fn new(rt: Runtime, service: Arc<SessionService>) -> Self {
        Self {
            rt,
            service,
            method: "vape".into(),
            amount: String::new(),
            unit: "puffs".into(),
            sessions: Vec::new(),
            status_message: String::new(),
        }
    }

    fn refresh_sessions(&mut self) {
        let service = self.service.clone();
        match self.rt.block_on(service.get_all_sessions()) {
            Ok(sessions) => self.sessions = sessions,
            Err(e) => self.status_message = format!("Error: {}", e),
        }
    }

    fn log_session(&mut self) {
        let service = self.service.clone();
        let method = self.method.clone();
        let amount = self.amount.clone();
        let unit = self.unit.clone();

        match self.rt.block_on(service.log_session(method, amount, unit)) {
            Ok(_) => {
                self.status_message = "Session logged!".into();
                self.amount.clear();
                self.refresh_sessions();
            }
            Err(e) => self.status_message = format!("Error: {}", e),
        }
    }
}

impl eframe::App for CannaNote {
    fn update(&mut self, ctx: &egui::Context, _frame: &mut eframe::Frame) {
        egui::CentralPanel::default().show(ctx, |ui| {
            ui.heading("CannaNote");
            ui.add_space(10.0);

            // Log session form
            ui.group(|ui| {
                ui.label("Log Session");

                ui.horizontal(|ui| {
                    ui.label("Method:");
                    egui::ComboBox::from_id_source("method")
                        .selected_text(&self.method)
                        .show_ui(ui, |ui| {
                            ui.selectable_value(&mut self.method, "vape".into(), "Vape");
                            ui.selectable_value(&mut self.method, "smoke".into(), "Smoke");
                            ui.selectable_value(&mut self.method, "edible".into(), "Edible");
                            ui.selectable_value(&mut self.method, "tincture".into(), "Tincture");
                        });
                });

                ui.horizontal(|ui| {
                    ui.label("Amount:");
                    ui.text_edit_singleline(&mut self.amount);
                    ui.label(&self.unit);
                });

                if ui.button("Log Session").clicked() {
                    self.log_session();
                }
            });

            ui.add_space(10.0);

            // Status message
            if !self.status_message.is_empty() {
                ui.label(&self.status_message);
            }

            ui.add_space(10.0);

            // Session list
            ui.group(|ui| {
                ui.horizontal(|ui| {
                    ui.label("Sessions");
                    if ui.button("Refresh").clicked() {
                        self.refresh_sessions();
                    }
                });

                ui.add_space(5.0);

                if self.sessions.is_empty() {
                    ui.label("No sessions yet. Log your first session above!");
                } else {
                    egui::ScrollArea::vertical()
                        .max_height(200.0)
                        .show(ui, |ui| {
                            for session in &self.sessions {
                                ui.horizontal(|ui| {
                                    ui.label(format!(
                                        "{} - {} {} ({})",
                                        session.timestamp.format("%H:%M"),
                                        session.amount,
                                        session.unit,
                                        session.method
                                    ));
                                });
                            }
                        });
                }
            });
        });
    }
}
