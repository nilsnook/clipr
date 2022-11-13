install:
	@echo "Installing clipr..."
	@go install ./cmd/clipr/

	@echo "Installing clipr daemon..."
	@go install ./cmd/cliprd/

	@echo "Preparing service for clipr daemon..."
	@./build/prod/prepservice

	@echo "Copying service file..."
	@mkdir -p ~/.config/systemd/user/
	@cp ./build/prod/cliprd.service ~/.config/systemd/user/

	@systemctl --user daemon-reload
	@echo "Enabling clipr daemon..."
	@systemctl --user enable cliprd.service
	@echo "Starting clipr daemon..."
	@systemctl --user start cliprd.service

	@echo "Clipr installed successfully!"

uninstall:
	@echo "Stopping clipr daemon..."
	@systemctl --user stop cliprd.service

	@echo "Disabling clipr daemon..."
	@systemctl --user disable cliprd.service

	@./build/prod/rmbins

	@echo "Cleaning config files..."
	@rm ~/.config/systemd/user/cliprd.service

	@echo "Clipr uninstalled successfully!"
